;redcode
;name RealImpKiller
;author Core War Community
;strategy Creates a continuous stream of DAT bombs
;assert 1

loop:   MOV #0, @3       ; Place DAT at pointer location
        ADD #1, 2        ; Increment pointer
        JMP -2           ; Continue
ptr:    DAT #0, #10      ; Pointer for bombing

end