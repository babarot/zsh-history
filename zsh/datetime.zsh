#!/usr/bin/env zsh

__zsh_history::datetime::today()
{
    local format="${1:-"+%F"}"
    date "$format"
}

__zsh_history::datetime::yesterday()
{
    # c.f. Can I get yesterday's date in a unix script?
    # https://community.hpe.com/t5/General/Can-I-get-yesterday-s-date-in-a-unix-script/td-p/3252466
    local format="${1:-"+%F"}"
    TZ="$TZ+24" date "$format"
}

__zsh_history::datetime::get_mtime()
{
    local file="${1:-$ZSH_HISTORY_FILE}"
    stat +mtime "$file"
}

__zsh_history::datetime::to_string()
{
    # e.g. __zsh_history::datetime::to_string 1486998000
    # => 2017-02-14
    strftime "%F" "${1:?}"
}

__zsh_history::datetime::to_epoch()
{
    # e.g. __zsh_history::datetime::to_epoch 2017-02-14
    # => 1486998000
    strftime -r "%F" "${1:?}"
}
