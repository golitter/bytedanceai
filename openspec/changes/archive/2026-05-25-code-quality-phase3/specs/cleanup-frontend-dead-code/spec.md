## ADDED Requirements

### Requirement: 移除未使用的 npm 依赖
系统 SHALL 从 `package.json` 中移除 `radix-ui` 包（从未被任何文件 import）。如果移除 `card.tsx`/`button.tsx`/`input.tsx` 后 `@radix-ui/react-slot` 和 `class-variance-authority` 不再被其他文件引用，则一并移除。

#### Scenario: 移除 radix-ui 后构建通过
- **WHEN** 从 package.json 移除 `radix-ui` 并执行 `pnpm install`
- **THEN** `pnpm build` 编译成功，无错误

#### Scenario: react-slot 和 class-variance-authority 随组件一并移除
- **WHEN** card.tsx、button.tsx、input.tsx 被删除且无其他文件使用这两个依赖
- **THEN** 从 package.json 移除 `@radix-ui/react-slot` 和 `class-variance-authority`，`pnpm build` 编译成功

### Requirement: 移除未使用的 shadcn/ui 组件文件
系统 SHALL 删除 `src/components/ui/card.tsx`、`src/components/ui/button.tsx`、`src/components/ui/input.tsx`，这三个组件从未被任何文件 import。

#### Scenario: 删除后构建通过
- **WHEN** 删除这三个文件
- **THEN** `pnpm build` 编译成功

### Requirement: 移除未使用的契约生成类型文件
系统 SHALL 删除 `src/generated/response.ts` 和 `src/generated/session.ts`，这两个文件导出的类型从未被任何文件 import。

#### Scenario: 删除后构建通过
- **WHEN** 删除这两个文件
- **THEN** `pnpm build` 编译成功

### Requirement: 移除 api.ts 中的死代码导出
系统 SHALL 从 `src/lib/api.ts` 中移除以下从未被外部使用的导出：`deleteTask()`、`patchSession()` 函数以及 `StreamEvent` 类型重导出（第 99 行）。

#### Scenario: 移除后无编译错误
- **WHEN** 移除这些导出
- **THEN** `pnpm build` 编译成功

### Requirement: 移除未使用的 Vite 脚手架资源
系统 SHALL 删除 `src/assets/hero.png`、`src/assets/react.svg`、`src/assets/vite.svg`，这些文件从未被任何源文件引用。

#### Scenario: 删除后构建通过
- **WHEN** 删除这三个资源文件
- **THEN** `pnpm build` 编译成功
