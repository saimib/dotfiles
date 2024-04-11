local status, _ = pcall(vim.cmd, "colorscheme nightfly")
if not status then
  print("colorsheme not found!")
  return
end
