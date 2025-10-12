# prompter

this is my prompter

## how to add as prompt
Run `make install` to install the binary. And make sure its on path:
```
which prompter
```
Then add a line as such in your bash file:
```
PROMPT='$(prompter)'
```

## Display
Then enjoy the custom prompt
```
[containerPrefix] [hostname] ~[pwd] 🔀<git-branch> <git-state-emoji><±num-of-changes>
⚡
```

It will try to reduce the length of the prompt if it is wider than the current
terminal window. 

