local status, autopairs = pcall(require, "nvim-autopairs")
if not status then
  print("nvim-autopairs not found!")
  return
end

autopairs.setup({
  check_ts = true,
  ts_config = {
    lua = { "string" },
    javascript = { "template_string"},
  },
})


local cmp_autopairs_status, cmp_autopairs = pcall(require, "nvim-autopairs.completion.cmp")
if not cmp_autopairs_status then
  print("nvim-autopairs.completion.cmp not found!")
  return
end

local cmp_status, cmp = pcall(require, "cmp")
if not cmp_status then
  print("cmp not found!")
  return
end

cmp.event:on("confirm_done", cmp_autopairs.on_confirm_done())
