mov r2, 5

loop:
    mov r0, 1
    mov r1, r2
    syscall           ; prints 5

    sub r2, 1         ; r2 becomes 4,3,2,1,0,-1,-2,...

    cmp r2, 0
    jz done           ; ← this should exit when r2==0

    jmp loop          ; ← this unconditional jump makes it infinite anyway!