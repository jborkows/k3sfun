-- Neovim configuration for Vikunja auto-transition Go project
-- This file configures native LSP (gopls) and golangci-lint to recognize
-- vikunja/auto-transition as the Go source root

-- Get the directory where .nvimrc.lua is located
-- Using debug info to get the current file's directory
local config_dir = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":p:h")

-- Set the Go root directory to vikunja/auto-transition
local go_root = vim.fs.joinpath(config_dir, "vikunja", "auto-transition")

-- Configure gopls using native LSP API (Neovim 0.10+)
vim.lsp.config("gopls", {
	cmd = { "gopls" },
	root_markers = { "go.mod" },
	root_dir = function(bufnr, callback)
		callback(go_root)
	end,
	settings = {
		gopls = {
			build = {
				-- Ensure we're working in the correct directory
				allowModfileModifications = false,
				allowImplicitNetworkAccess = false,
			},
			ui = {
				diagnostic = {
					-- Enable diagnostics
					annotations = {
						bounds = true,
						escape = true,
						inline = true,
						["nil"] = true,
					},
				},
			},
			formatting = {
				gofumpt = true,
			},
			-- Enable additional analyses
			analyses = {
				unusedparams = true,
				shadow = true,
				nilness = true,
				unusedwrite = true,
				useany = true,
			},
			staticcheck = true,
		},
	},
})

-- Enable gopls for Go files
vim.lsp.enable("gopls")

-- Configure golangci-lint using nvim-lint (if available)
local has_lint, lint = pcall(require, "lint")
if has_lint then
	-- Set up golangci-lint for Go files
	lint.linters_by_ft = lint.linters_by_ft or {}
	lint.linters_by_ft.go = { "golangcilint" }

	-- Configure golangci-lint for the auto-transition project
	local golangcilint = require("lint.linters.golangcilint")
	golangcilint.cmd = "golangci-lint"
	golangcilint.args = {
		"run",
		"--no-config",
		"--output.json.path", "stdout",
		"--issues-exit-code", "0",
		"--show-stats=false",
	}
	golangcilint.cwd = go_root

	-- Auto-run linter on save for Go files
	vim.api.nvim_create_autocmd({ "BufWritePost" }, {
		pattern = "*/vikunja/auto-transition/*.go",
		callback = function()
			lint.try_lint()
		end,
	})

	vim.notify("golangci-lint configured for Go files", vim.log.levels.INFO)
else
	vim.notify("nvim-lint not found - golangci-lint not configured", vim.log.levels.WARN)
end

-- Alternative: Run golangci-lint manually via command
vim.api.nvim_create_user_command("GolangciLint", function()
	local bufnr = vim.api.nvim_get_current_buf()
	local filename = vim.api.nvim_buf_get_name(bufnr)

	-- Run golangci-lint and capture output
	vim.fn.jobstart({
		"golangci-lint",
		"run",
		"--out-format", "line-number",
		filename,
	}, {
		cwd = go_root,
		on_stdout = function(_, data)
			if data then
				for _, line in ipairs(data) do
					if line ~= "" then
						vim.notify(line, vim.log.levels.WARN)
					end
				end
				end
		end,
		on_stderr = function(_, data)
			if data then
				for _, line in ipairs(data) do
					if line ~= "" then
						vim.notify(line, vim.log.levels.ERROR)
					end
				end
				end
		end,
	})
end, {
	desc = "Run golangci-lint on current file",
})

-- Set working directory to the Go project root when opening Go files
vim.api.nvim_create_autocmd("BufEnter", {
	pattern = "*/vikunja/auto-transition/*.go",
	callback = function()
		vim.cmd("lcd " .. go_root)
	end,
})

-- Set up proper Go module path handling
vim.env.GO111MODULE = "on"

-- Print a message to confirm the configuration is loaded
vim.notify("Go LSP configured for: " .. go_root, vim.log.levels.INFO)
