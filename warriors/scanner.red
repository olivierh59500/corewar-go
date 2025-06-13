;redcode
;name Scanner
;author Core War Community
;strategy Scans for enemy code and bombs it
;assert 1

        ADD #10, ptr     ; Increment pointer
scan:   CMP @ptr, #0     ; Compare location with zero
        JMP found        ; If not zero, enemy found
        ADD #1, ptr      ; Increment scan pointer
        JMP scan         ; Continue scanning
found:  MOV #0, @ptr     ; Bomb the location
        JMP -5           ; Go back to main loop
ptr:    DAT #0, #0       ; Pointer

end