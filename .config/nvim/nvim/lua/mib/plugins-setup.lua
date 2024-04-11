local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
vim.fn.system({
  "git",
  "clone",
  "--filter=blob:none",
  "https://github.com/folke/lazy.nvim.git",
"--branch=stable", -- latest stable release
  lazypath,
})
end

vim.opt.rtp:prepend(lazypath)

vim.g.mapleader = " "
vim.g.maplocalleader = "\\"

require("lazy").setup({
  {"bluz71/vim-nightfly-colors", name="nightfly", lazy=false, priority=1000},
  "nvim-lua/plenary.nvim",
  "christoomey/vim-tmux-navigator",
  "szw/vim-maximizer",
  "tpope/vim-surround",
  "vim-scripts/ReplaceWithRegister",
  {"numToStr/Comment.nvim", opts={}, lazy=false},
  {"nvim-tree/nvim-tree.lua", version="*", lazy=false, dependencies = {"nvim-tree/nvim-web-devicons"},
  config = function()
    require("nvim-tree").setup()
  end},
  {"nvim-lualine/lualine.nvim", dependencies={"nvim-tree/nvim-web-devicons"}, 
  config = function()
    require("lualine").setup({options = {theme = "gruvbox"}})
  end},
  {"nvim-telescope/telescope-fzf-native.nvim", build = "make" },
  {"nvim-telescope/telescope.nvim", branch = "0.1.x", dependencies = { "nvim-lua/plenary.nvim"}},
  "hrsh7th/nvim-cmp",
  "hrsh7th/cmp-buffer",
  "hrsh7th/cmp-path",
  "L3MON4D3/LuaSnip",
  "saadparwaiz1/cmp_luasnip",
  "rafamadriz/friendly-snippets",
  "williamboman/mason.nvim",
  "williamboman/mason-lspconfig.nvim",
  "neovim/nvim-lspconfig",
  "hrsh7th/cmp-nvim-lsp",
  {"glepnir/lspsaga.nvim", branch = "main" },
  "onsails/lspkind.nvim",
  "nvimtools/none-ls.nvim",
  "jayp0521/mason-null-ls.nvim",
  {"nvim-treesitter/nvim-treesitter", build = ":TSUpdate" },
  "windwp/nvim-autopairs",
  "windwp/nvim-ts-autotag",
  {"lewis6991/gitsigns.nvim",
  config = function()
    require("gitsigns").setup()
  end}
})
