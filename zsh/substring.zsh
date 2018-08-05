#!/usr/bin/env zsh

__zsh_history::substring::search_begin()
{
    setopt localoptions extendedglob

    _history_substring_search_refresh_display=0
    _history_substring_search_query_highlight=

    if [[ -z $BUFFER || $BUFFER != $_history_substring_search_result ]]; then
        _history_substring_search_query=$BUFFER
        _history_substring_search_query_escaped=${BUFFER//(#m)[\][()|\\*?#<>~^]/\\$MATCH}

        local   select="SELECT command FROM history"
        local    where="WHERE dir = '$PWD' AND command LIKE '$_history_substring_search_query_escaped%'"
        local group_by="GROUP BY command ORDER BY id ASC"
        _history_substring_search_matches=(${(@f)"$(zhist -q "$select $where $group_by")"})
        if [[ $#_history_substring_search_matches -eq 0 ]]; then
            _history_substring_search_matches=()
        fi

        _history_substring_search_matches_count=$#_history_substring_search_matches

        if [[ $WIDGET == history-substring-search-up ]]; then
            _history_substring_search_match_index=$(( _history_substring_search_matches_count + 1 ))
        else
            _history_substring_search_match_index=$_history_substring_search_matches_count
        fi
    fi
}

__zsh_history::substring::search_end()
{
    setopt localoptions extendedglob

    _history_substring_search_result=$BUFFER

    if (( $_history_substring_search_refresh_display == 1 )); then
        region_highlight=()
        CURSOR=$#BUFFER
    fi

    # highlight command line using zsh-syntax-highlighting
    if (( $+functions[_zsh_highlight] )); then
        _zsh_highlight
    fi

    # highlight the search query inside the command line
    if [[ -n $_history_substring_search_query_highlight && -n $_history_substring_search_query ]]; then
        : ${(S)BUFFER##(#m$HISTORY_SUBSTRING_SEARCH_GLOBBING_FLAGS)($_history_substring_search_query##)}
        local begin=$(( MBEGIN - 1 ))
        local end=$(( begin + $#_history_substring_search_query ))
        region_highlight+=("$begin $end $_history_substring_search_query_highlight")
    fi
}

__zsh_history::substring::history_up()
{
    _history_substring_search_refresh_display=1

    if (( $_history_substring_search_match_index > 0 )); then
        BUFFER=$_history_substring_search_matches[$_history_substring_search_match_index]
        _history_substring_search_query_highlight=$HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND
        (( _history_substring_search_match_index-- ))
    else
        __zsh_history::substring::not_found
    fi
}

__zsh_history::substring::history_down()
{
    _history_substring_search_refresh_display=1

    if (( _history_substring_search_match_index < $#_history_substring_search_matches )); then
        (( _history_substring_search_match_index++ ))
        BUFFER=$_history_substring_search_matches[$_history_substring_search_match_index]
        _history_substring_search_query_highlight=$HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND
    else
        BUFFER=$_history_substring_search_old_buffer
        _history_substring_search_query_highlight=
    fi
}

__zsh_history::substring::not_found()
{
    _history_substring_search_old_buffer=$BUFFER
    BUFFER=$_history_substring_search_query
    _history_substring_search_query_highlight=$HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_NOT_FOUND
}

__zsh_history::substring::reset()
{
    _history_substring_search_result=
}
