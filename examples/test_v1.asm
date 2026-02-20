; leg-lock-lang example program - tests most implemented features
; Features demonstrated:
;   - labels & jumps (forward + backward)
;   - mov reg, reg
;   - mov reg, immediate (positive + negative)
;   - add reg, reg
;   - add reg, immediate
;   - sub reg, reg
;   - sub reg, immediate
;   - syscall print (r0=1, r1=value)
;   - syscall exit  (r0=0, r1=exitcode)
;   - comments
;   - multiple instructions per line? (no - one per line currently)

start:
    mov r0,  42          ; load constant
    mov r1,  r0          ; reg → reg copy
    add r1,  8           ; reg += immediate
    add r2,  r1          ; r2 = r1 (because r2 was 0)
    sub r2,  10          ; r2 -= 10
    mov r3,  -5          ; negative immediate

    ; print 50  (42 + 8)
    mov r0,  1
    mov r1,  r2
    syscall              ; should print 50

    ; print -5
    mov r0,  1
    mov r1,  r3
    syscall              ; should print -5

    ; simple loop (count down from 5 → 1)
    mov r4,  5
loop:
    mov r0,  1
    mov r1,  r4
    syscall              ; print 5 4 3 2 1

    sub r4,  1
    mov r5,  r4
    add r5,  1           ; r5 = r4 + 1 (just for testing add reg,reg)

    ; if r4 > 0 → jump back
    mov r6,  r4
    sub r6,  0           ; dummy to set flags? (we don't have cmp/branch yet)
    ; (so we fake condition with jump always for now - real condition later)
    jmp loop             ; WARNING: infinite loop right now!

    ; We never reach here with current code (no conditional branch)
    ; This is just to show unreachable code after infinite loop

    mov r0,  1
    mov r1,  999         ; should NOT be printed
    syscall

done:
    mov r0,  0           ; exit syscall
    mov r1,  7           ; exit code 7
    syscall