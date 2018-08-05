#!/usr/bin/env zsh

zmodload zsh/stat
zmodload zsh/datetime

export ZSH_HISTORY_FILE=${ZSH_HISTORY_FILE:-"$HOME/.zsh_history.db"}
ZSH_HISTORY_BACKUP_DIR=${ZSH_HISTORY_BACKUP_DIR:-"$HOME/.zsh/history/backup"}
ZSH_HISTORY_FILTER=${ZSH_HISTORY_FILTER:-"fzy:fzf-tmux:fzf:peco"}
ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND="bg=magenta,fg=white,bold"
ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_NOT_FOUND="bg=red,fg=white,bold"
ZSH_HISTORY_SUBSTRING_SEARCH_GLOBBING_FLAGS="i"

#
# Keybindings
#

if [[ -n $ZSH_HISTORY_KEYBIND_GET_BY_DIR ]]; then
    zle -N "__zsh_history::keybind::get_by_dir"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_BY_DIR" "__zsh_history::keybind::get_by_dir"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_GET_ALL ]]; then
    zle -N "__zsh_history::keybind::get_all"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_ALL" "__zsh_history::keybind::get_all"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_SCREEN ]]; then
    zle -N "__zsh_history::keybind::screen"
    bindkey "$ZSH_HISTORY_KEYBIND_SCREEN" "__zsh_history::keybind::screen"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_ARROW_UP ]]; then
    zle -N "__zsh_history::keybind::arrow_up"
    bindkey "$ZSH_HISTORY_KEYBIND_ARROW_UP" "__zsh_history::keybind::arrow_up"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_ARROW_DOWN ]]; then
    zle -N "__zsh_history::keybind::arrow_down"
    bindkey "$ZSH_HISTORY_KEYBIND_ARROW_DOWN" "__zsh_history::keybind::arrow_down"
fi

#
# Configurations
#

if [[ $ZSH_HISTORY_CASE_SENSITIVE == true ]]; then
    unset ZSH_HISTORY_SUBSTRING_SEARCH_GLOBBING_FLAGS
fi

if [[ $ZSH_HISTORY_DISABLE_COLOR == true ]]; then
    unset ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND
    unset ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_NOT_FOUND
fi

#
# Loading
#

for f in "${0:A:h}"/zsh/*.zsh(N-.)
do
    source "$f"
done
unset f
