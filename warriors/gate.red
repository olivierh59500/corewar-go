;redcode
;name Gate
;author Core War Community  
;strategy Scanner that creates a gate between two locations
;assert 1

        CMP 9, 19       ; Compare two locations
        JMP found       ; If different, something found
        ADD #1, -1      ; Increment scan
        JMP -3          ; Loop back
found:  SPL 0           ; Split when enemy found
        MOV @-5, <-4    ; Create gate
        ADD #1, -1      ; Increment
        JMP -2          ; Loop

end