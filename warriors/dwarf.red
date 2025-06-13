;redcode
;name Dwarf
;author A.K. Dewdney
;strategy Bombs core with DAT instructions
;assert 1

        ADD #4, 3        ; Add 4 to the B field of the DAT instruction
        MOV #0, @2       ; Move 0 to the location pointed by the DAT
        JMP -2           ; Jump back to ADD
        DAT #0, #0       ; This serves as the pointer

end