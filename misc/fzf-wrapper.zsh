#!/bin/zsh

local -i ret=0
local -a files

if [[ -z $argv[1] ]]; then
    echo "too few arguments" >&2
    exit 1
fi

if (( ! $+commands[zhist] )); then
    echo "zhist: command not found" >&2
    exit 1
fi

IFS=$'\n'

files=( ${(@f):-"$(fzf)"} )

if (( $#files == 0 )); then
    # No files selected
    exit 0
fi

${=argv[@]} ${=files[@]}
ret=$status

# Insert selected files as one of the zsh history
zhist -i "${(j: :)argv[@]} ${(j: :)files[@]}" "$ret" \
    &>/dev/null

exit $ret
