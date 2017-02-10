#!/bin/zsh

export ZSH_HISTORY_FILE=${ZSH_HISTORY_FILE:-"$HOME/.zsh_history.db"}
export ZSH_HISTORY_FILTER=${ZSH_HISTORY_FILTER:-"fzy:fzf-tmux:fzf:peco"}

if [[ -n $ZSH_HISTORY_KEYBIND_GET_BY_DIR ]]; then
    zle -N "__zsh_history::keybind::get_by_dir"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_BY_DIR" "__zsh_history::keybind::get_by_dir"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_GET_ALL ]]; then
    zle -N "__zsh_history::keybind::get_all"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_ALL" "__zsh_history::keybind::get_all"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_INTERACTIVE ]]; then
    zle -N "__zsh_history::keybind::interactive"
    bindkey "$ZSH_HISTORY_KEYBIND_INTERACTIVE" "__zsh_history::keybind::interactive"
fi

for f in "${0:A:h}"/zsh/*.zsh(N-.)
do
    source "$f"
done
unset f
