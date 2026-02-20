    mov r4, 5           ; counter

loop:
    mov r0, 1
    mov r1, r4
    syscall             ; print 5, then 4, ...

    sub r4, 1

    mov r6, r4          ; copy counter to r6
    add r6, 1           ; r6 = r4 + 1 (just demo)

    ; Very crude "if r4 == 0 then skip jump" simulation
    ; (but without real compare → we can't do it properly yet)

    ; Temporary hack: repeat jmp only a few times manually
    ; or just let it run once and remove jmp

    ; jmp loop          ; ← comment this for finite run