set nocompatible
filetype off

set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()

Plugin 'gmarik/Vundle.vim'
Plugin 'ogier/guessindent'
Plugin 'rking/ag.vim'
Plugin 'metacosm.vim'
"Plugin 'fatih/vim-go'
"Plugin 'fatih/molokai'
Plugin 'majutsushi/tagbar'
"Plugin 'Valloric/YouCompleteMe'
"Plugin 'clang-complete'
"Plugin 'vim-clang'
"Plugin 'vim-airline/vim-airline'
"Plugin 'vim-airline/vim-airline-themes'

call vundle#end()

filetype plugin indent on

" spellcheck
setlocal nospell spelllang=en_us
set spellfile=~/.vim/spell/words.utf-8.add
autocmd FileType gitcommit setlocal spell
autocmd BufReadPost .stgit-edit.txt setlocal spell
nnoremap <leader>spell :setlocal spell<cr>
nnoremap <leader>nospell :setlocal nospell<cr>

" completion
set complete+=kspell
set completeopt+=longest
set wildmenu	" use the tab completion menu
set showcmd	" show us the command we're typing

" color
syntax on
colorscheme metacosm

"let g:airline#extensions#tabline#enabled = 1
"let g:airline_theme='luna'
"set statusline=%F%m%r%h%w[%L][%{&ff}]%y[%p%%][%04l,%04v]
"set laststatus=2
"set noshowmode

" wrapping
"set textwidth=80
"set wrapmargin=80
set nowrap

" indentation
set autoindent
set smartindent
set cindent
set smarttab
set noexpandtab
set tabstop=8
set shiftwidth=8
set softtabstop=8
" set cinoptions=(0,u0,U0,l1

" bad formatting
highlight BadFormat ctermbg=darkred guibg=darkred

func UpdateBadFormat()
	if &expandtab
		match BadFormat /\s\+$\|\t\|\%81v/
	elseif &tabstop == 8
		match BadFormat /\s\+$\|\ \+\t\| \{8,}\|\%81v/
	else
		match BadFormat /\s\+$\|\ \+\t\| \{4,}\|\%81v/
	endif
endfun

let g:guessindent_prefer_tabs = 1
autocmd BufReadPost * :GuessIndent
autocmd BufEnter * call UpdateBadFormat()

" git grep
func GitGrep(...)
  let save = &grepprg
  set grepprg=git\ grep\ -n\ $*
  let s = 'grep'
  for i in a:000
    let s = s . ' ' . i
  endfor
  exe s
  let &grepprg = save
endfun

func GitGrepWord()
  normal! "zyiw
  call GitGrep('-w -e ', getreg('z'))
endf

command -nargs=? G call GitGrep(<f-args>)

" folding
set nofoldenable
set foldlevelstart=99

" interface
behave mswin

set backspace=indent,eol,start " enable backspacing over everything
set whichwrap=b,s,<,>,[,] " wrap over lines

set ruler	" show the cursor position at all times
set number	" show line numbers
set showmatch	" show bracket pair

set title	" xterm title
set noerrorbells	" no beeps on error messages
set novisualbell
set belloff=all
set mouse=a	" enable mouse usage (all modes) in terminals
"set guifont=Terminus\ 14
"set guifont=SourceSansPro\ 14
set guifont=InconsolataLGC\ 14
set guioptions-=T	"remove toolbar
set guioptions-=m	"remove menu bar
set guioptions-=r	"remove right scrollbar
set guioptions-=L	"remove left scrollbar
set scrolloff=3		"minimal number of screen lines to keep above and below the cursor
cwindow		" show an error window when there are errors

set lazyredraw		"don't redraw screen while executing macros/mappings
set ttyfast

" search
set incsearch
set hlsearch
set ignorecase
set smartcase

set history=1000
set undolevels=1000

set writeany
set hidden " allow movement to another buffer without saving the current one
set autowriteall
set autoread " read changes from disk automatically
set confirm

"set backup
"set backupdir=$HOME/.vim/backup
"set directory=$HOME/.vim/tmp

set encoding=utf-8
set fileencodings=utf-8

set tags=tags;/,~/src/common-tags,~/src/linux/tags

" HOTKEYS

noremap <C-s> :update<cr>
inoremap <C-s> <C-o>:update<cr>
vnoremap <C-s> <C-c>:update<cr>

noremap <C-z> u
inoremap <C-z> <C-o>u

inoremap <C-r> <C-o><C-r>

vnoremap <C-c> "+y
vnoremap <C-x> "+x

noremap <C-q> <C-v>

map <C-v> "+gP
inoremap <C-v> <C-o>"+gP
cmap <C-v> <C-r>+

noremap <C-a> gggH<C-o>G
inoremap <C-a> <C-o>gg<C-o>gH<C-o>G
cnoremap <C-a> <C-c>gggH<C-o>G
onoremap <C-a> <C-c>gggH<C-o>G
snoremap <C-a> <C-c>gggH<C-o>G
xnoremap <C-a> <C-c>ggVG

vnoremap <BS> d

nmap <F8> :TagbarToggle<cr>
imap <F8> <C-o>:TagbarToggle<cr>

nmap <F9> :make!<cr>
imap <F9> <C-o>:make!<cr>

nmap <M-right> :call GitGrepWord()<cr>
imap <M-right> <C-o>:call GitGrepWord()<cr>

nmap <M-down> :cn<cr>
imap <M-down> <C-o>:cn<cr>

nmap <M-up> :cp<cr>
imap <M-up> <C-o>:cp<cr>

nmap <C-right> <C-]>
imap <C-right> <C-o><C-]>

nmap <c-right> <C-]>
imap <C-right> <C-o><C-]>

nmap <C-left> <C-t>
imap <C-left> <C-o><C-t>

nmap <C-down> :tn<cr>
imap <C-down> <C-o>:tn<cr>

nmap <C-up> :tp<cr>
imap <C-up> <C-o>:tp<cr>

nmap <C-tab> :bn<cr>
imap <C-tab> <C-o>:bn<cr>

nmap <C-S-tab> :bp<cr>
imap <C-S-tab> <C-o>:bp<cr>
