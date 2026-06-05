package impl

import (
	"fmt"
	"os"
	"strings"
	"time"

	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/model"
	"agenthub/backend/internal/service"
	"agenthub/backend/pkg/agentend_client"
)

type SkillService struct {
	skillDao    dao.SkillDao
	sessionDao  dao.SessionDao
	agentClient *agentend_client.Client
}

func NewSkillService(skillDao dao.SkillDao, sessionDao dao.SessionDao, agentClient *agentend_client.Client) *SkillService {
	return &SkillService{
		skillDao:    skillDao,
		sessionDao:  sessionDao,
		agentClient: agentClient,
	}
}

func (svc *SkillService) UploadSkill(filename string, zipData []byte) (*service.ValidationResult, error) {
	result, tmpDir, err := service.ValidateZip(zipData)
	if err != nil {
		return nil, service.ErrInternal("invalid zip file")
	}

	if result.Valid {
		zipName := strings.TrimSuffix(filename, ".zip")
		if zipName != result.Name {
			_ = os.RemoveAll(tmpDir)
			return nil, service.ErrBadRequest(fmt.Sprintf("zip filename (%s) must match SKILL.md name (%s)", zipName, result.Name))
		}

		count, err := svc.skillDao.CountBuiltinByName(result.Name)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			_ = os.RemoveAll(tmpDir)
			return &service.ValidationResult{
				Valid:  false,
				Errors: []string{"name conflicts with builtin skill"},
			}, nil
		}
	}

	return result, nil
}

func (svc *SkillService) ConfirmSkill(name, description string, fileCount int, totalSize int64, tmpDir string) (*service.SkillImportResult, error) {
	count, err := svc.skillDao.CountByName(name)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, service.ErrConflict("skill name already exists")
	}

	zipData, err := service.PackValidatedSkillDir(name, tmpDir)
	if err != nil {
		return nil, err
	}

	if err := svc.skillDao.CreateSkill(model.SkillHub{
		Name:        name,
		Builtin:     false,
		Description: description,
		FileCount:   fileCount,
		TotalSize:   totalSize,
		Content:     zipData,
	}); err != nil {
		return nil, err
	}

	_ = os.RemoveAll(tmpDir)
	return &service.SkillImportResult{Success: true, Name: name}, nil
}

func (svc *SkillService) ListSkills() ([]service.SkillHubItem, error) {
	skills, err := svc.skillDao.ListSkills()
	if err != nil {
		return nil, err
	}

	items := make([]service.SkillHubItem, 0, len(skills))
	for _, skill := range skills {
		importCount, err := svc.skillDao.CountImportsBySkillName(skill.Name)
		if err != nil {
			return nil, err
		}
		items = append(items, service.SkillHubItem{
			Name:        skill.Name,
			Builtin:     skill.Builtin,
			Description: skill.Description,
			FileCount:   skill.FileCount,
			TotalSize:   skill.TotalSize,
			ImportCount: importCount,
			CreatedAt:   skill.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items, nil
}

func (svc *SkillService) DeleteSkill(name string) error {
	skill, err := svc.skillDao.GetSkillByName(name)
	if err != nil {
		return err
	}
	if skill == nil {
		return service.ErrNotFound("skill not found")
	}
	if skill.Builtin {
		return service.ErrForbidden("cannot delete builtin skill")
	}
	return svc.skillDao.DeleteSkillCascade(name)
}

func (svc *SkillService) ImportSkill(skillName, sessionID string) (*service.SkillImportResult, error) {
	session, err := svc.sessionDao.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, service.ErrNotFound("session not found")
	}

	allowedTypes := map[string]bool{"claude-code": true, "opencode": true, "codex": true}
	if !allowedTypes[session.AgentType] {
		return nil, service.ErrForbidden("orchestrator does not support importing external skills")
	}

	skill, err := svc.skillDao.GetSkillByName(skillName)
	if err != nil {
		return nil, err
	}
	if skill == nil {
		return nil, service.ErrNotFound("skill not found in hub")
	}

	exists, err := svc.skillDao.HasAgentSkill(sessionID, skillName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, service.ErrConflict("skill already imported to this session")
	}

	zipData, err := svc.skillDao.GetSkillContent(skillName)
	if err != nil {
		return nil, err
	}
	if len(zipData) == 0 {
		return nil, service.ErrInternal("pack skill files failed: no zip data")
	}

	if err := svc.agentClient.InstallSkill(session.AgentType, sessionID, skillName, zipData); err != nil {
		return nil, service.ErrInternal("install skill to worktree failed: " + err.Error())
	}

	if err := svc.skillDao.CreateAgentSkill(model.AgentSkill{
		SessionID:  sessionID,
		SkillName:  skillName,
		AgentType:  session.AgentType,
		ImportedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return &service.SkillImportResult{
		Success: true,
		Skill:   skillName,
		Session: sessionID,
	}, nil
}

func (svc *SkillService) RemoveSkill(skillName, sessionID string) (*service.SkillImportResult, error) {
	session, err := svc.sessionDao.GetBySessionID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, service.ErrNotFound("session not found")
	}

	if err := svc.agentClient.RemoveSkill(session.AgentType, sessionID, skillName); err != nil {
		return nil, service.ErrInternal("remove skill files from worktree failed: " + err.Error())
	}

	if err := svc.skillDao.DeleteAgentSkill(sessionID, skillName); err != nil {
		return nil, err
	}

	return &service.SkillImportResult{
		Success: true,
		Skill:   skillName,
		Session: sessionID,
	}, nil
}

func (svc *SkillService) ReportBuiltinSkills(skills []service.BuiltinSkillItem) error {
	for _, skill := range skills {
		if err := svc.skillDao.UpsertSkillHub(skill.Name, skill.Description, true); err != nil {
			return err
		}
	}
	return nil
}
