#!/bin/zsh

__zsh_history::keybind::get_all()
{
    BUFFER="$(
    __zsh_history::history::get \
        "SELECT DISTINCT(command) FROM history ORDER BY id DESC" \
        "$LBUFFER"
    )"
    CURSOR=$#BUFFER
    zle reset-prompt
}

__zsh_history::keybind::get_by_dir()
{
    BUFFER="$(
    __zsh_history::history::get \
        "SELECT DISTINCT(command) FROM history WHERE dir = '$PWD' ORDER BY id DESC" \
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

    # Launch with screen
    local res="$(zhist -s $LBUFFER)"
    if [[ -n $res ]]; then
        BUFFER="$res"
        CURSOR=$#BUFFER
    fi

    zle reset-prompt
}

__zsh_history::keybind::arrow_up()
{
    __zsh_history::substring::search_begin
    __zsh_history::substring::history_up
    __zsh_history::substring::search_end
}

__zsh_history::keybind::arrow_down()
{
    __zsh_history::substring::search_begin
    __zsh_history::substring::history_down
    __zsh_history::substring::search_end
}
