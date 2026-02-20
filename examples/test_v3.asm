; Init memory
mov r0, 100     ; address
mov [100], 42   ; store 42 at addr 100

; Load from mem
mov r1, [100]   ; r1 = 42

; Compare and conditional jump
cmp r1, 42
jne error       ; shouldn't jump

mov r2, 5       ; loop counter

loop:
    cmp r2, 0
    jz done     ; exit if zero

    mov r0, 1
    mov r1, r2
    syscall     ; print r2 (5 4 3 2 1)

    sub r2, 1
    jmp loop

done:
    mov r0, 0
    mov r1, 0
    syscall     ; exit 0

error:
    mov r0, 0
    mov r1, 1
    syscall     ; exit 1