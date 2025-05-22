# prompter

this is my prompter

## how to add as prompt
Run `make install` to install the binary. And make sure its on path:
```
which prompter
```
Then add a line as such in your zsh file:
```
PROMPT='$(prompter)'
```

## Display
Then enjoy the custom prompt
```
~[pwd] 🔀<git-branch> <git-state-emoji><±num-of-changes>
⚡
```

