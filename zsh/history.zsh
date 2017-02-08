#!/bin/zsh

__zsh_history::history::add()
{
    local status_code="$status"
    local last_command="$(fc -ln -1)"

    if (( ! $+commands[zhist] )); then
        return 1
    fi

    zhist -a "$last_command" "$status_code"
}

__zsh_history::history::get()
{
    local filter query="${1:?}" lb="$2"
    filter="$(__zsh_history::utils::get_filter "$ZSH_HISTORY_FILTER")"

    if [[ -z $filter ]]; then
        print -r -- >&2 'zsh-history-enhanced: ZSH_HISTORY_FILTER is an invalid'
        return 1
    fi

    if (( ! $+commands[zhist] )); then
        return 1
    fi

    zhist -q "$query" \
        | __zsh_history::filter::grep "$lb" \
        | ${=filter} \
        | __zsh_history::filter::remove_ansi
}
