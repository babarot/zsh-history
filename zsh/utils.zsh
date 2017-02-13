#!/bin/zsh

__zsh_history::utils::get_filter()
{
    local x candidates

    if [[ -z $1 ]]; then
        return 1
    fi

    # candidates should be list like "a:b:c" concatenated by a colon
    candidates="$1:"

    while [[ -n $candidates ]]
    do
        # the first remaining entry
        x=${candidates%%:*}
        # reset candidates
        candidates=${candidates#*:}

        if type "${x%% *}" &>/dev/null; then
            echo "$x"
            return 0
        else
            continue
        fi
    done

    return 1
}

__zsh_history::utils::backup_db()
{
    local hist_dir="$ZSH_HISTORY_BACKUP_DIR/$(date +%Y/%m)"
    local backup_db="$hist_dir/$(__zsh_history::datetime::yesterday "+%d").db"
    local mtime ymd_mtime epoch_mtime
    local ymd_yesterday epoch_yesterday

    if [[ -f $backup_db ]]; then
        return 0
    else
        mkdir -p "$hist_dir"
    fi

    mtime="$(__zsh_history::datetime::get_mtime)"
    ymd_mtime="$(__zsh_history::datetime::to_string "$mtime")"
    epoch_mtime="$(__zsh_history::datetime::to_epoch "$ymd_mtime")"
    ymd_yesterday="$(__zsh_history::datetime::today)" # today
    epoch_yesterday="$(__zsh_history::datetime::to_epoch "$ymd_yesterday")"

    if (( epoch_mtime >= epoch_yesterday )); then
        command cp -f "$ZSH_HISTORY_FILE" "$backup_db"
    fi
}
