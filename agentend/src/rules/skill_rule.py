from src.rules.base import BaseRule


class SkillRule(BaseRule):
    name = "skill"
    description = "Injects output skill prompt"
    phase = "pre"
    priority = 1

    def check(self, context: dict) -> bool:
        return True

    def enforce(self, context: dict) -> dict:
        return {
            "system_prompt_append": (
                "## 输出技能\n"
                "\n"
                "workspace 中有 `render` 工具（路径在 skills 目录下），提供 5 个子命令：\n"
                "- `./render html-render '<div>...</div>'` — 渲染 HTML（也支持 stdin）\n"
                "- `./render image <path>` — 展示工作区图片\n"
                "- `./render attachment <path>` — 提供文件下载\n"
                "- `./render diff` — 展示工作区文件变更\n"
                "- `./render preview <url>` — 展示网页预览\n"
                "\n"
                "当你需要输出富媒体内容时，必须调用 `render` 工具，将它的 stdout 原样包含在你的回复中。\n"
                "不要手动编写 aka_yhy 格式，必须通过调用 render 工具生成。\n"
                "普通文本直接输出即可，不需要调用工具。"
            ),
        }
