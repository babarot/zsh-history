#!/bin/zsh

autoload -Uz add-zsh-hook

add-zsh-hook precmd "__zsh_history::utils::backup_db"

add-zsh-hook precmd "__zsh_history::history::add"

add-zsh-hook preexec "__zsh_history::substring::reset"
