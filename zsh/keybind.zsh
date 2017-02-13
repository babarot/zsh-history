#!/bin/zsh

__zsh_history::keybind::get_all()
{
    BUFFER="$(
    __zsh_history::history::get \
        "select distinct(command) from history order by id desc" \
        "$LBUFFER"
    )"
    CURSOR=$#BUFFER
    zle reset-prompt
}

__zsh_history::keybind::get_by_dir()
{
    BUFFER="$(
    __zsh_history::history::get \
        "select distinct(command) from history where dir = '$PWD' and status == 0 order by id desc" \
        "$LBUFFER"
    )"
    CURSOR=$#BUFFER
    zle reset-prompt
}

__zsh_history::keybind::screen()
{
    if (( ! $+commands[zhist] )); then
        return 1
    fi

    local res="$(zhist -s $LBUFFER)"
    if [[ -n $res ]]; then
        BUFFER="$res"
        CURSOR=$#BUFFER
    fi

    zle reset-prompt
}
