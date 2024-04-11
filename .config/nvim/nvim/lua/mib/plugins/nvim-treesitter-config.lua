local status, treesitter_configs = pcall(require, "nvim-treesitter.configs")
if not status then
  print("treesitter not found!")
  return
end

treesitter_configs.setup({
  ensure_installed = { "go", "javascript", "lua", "vim"},
  sync_install = false,
  highlight = { enable = true },
  indent = { enable = true },
  autotag = { enable = true },
  auto_install = true,
})
