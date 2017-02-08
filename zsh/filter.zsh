#!/bin/zsh

__zsh_history::filter::grep()
{
    if [[ -z $1 ]]; then
        cat -
    else
        grep --color="always" "$1"
    fi
}

__zsh_history::filter::remove_ansi()
{
    perl -pe 's/\e\[?.*?[\@-~]//g'
}
