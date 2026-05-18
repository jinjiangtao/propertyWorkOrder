try {
    Write-Host "Test 1: Register user"
    $resp = Invoke-RestMethod "http://localhost:8080/api/register" -Method Post -ContentType "application/json" -Body '{"username":"testuser","password":"123456"}'
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`nTest 2: Login user"
    $resp = Invoke-RestMethod "http://localhost:8080/api/login" -Method Post -ContentType "application/json" -Body '{"username":"testuser","password":"123456"}'
    Write-Host "Success: $($resp | ConvertTo-Json)"
    $userId = $resp.user_id

    Write-Host "`nTest 3: Create repair"
    $body = '{"user_id":' + $userId + ',"username":"testuser","repair_type":"水电维修","description":"水龙头漏水"}'
    $resp = Invoke-RestMethod "http://localhost:8080/api/repair/create" -Method Post -ContentType "application/json" -Body $body
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`nTest 4: Admin get all repairs"
    $resp = Invoke-RestMethod "http://localhost:8080/api/repair/all" -Method Get
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`nTest 5: Create worker"
    $resp = Invoke-RestMethod "http://localhost:8080/api/worker/create" -Method Post -ContentType "application/json" -Body '{"work_no":"W001","name":"Worker1","phone":"13800138001","password":"123456","skill_type":"水电"}'
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`nTest 6: Assign worker to repair"
    $resp = Invoke-RestMethod "http://localhost:8080/api/repair/assign" -Method Post -ContentType "application/json" -Body '{"repair_id":1,"worker_id":1}'
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`nTest 7: Worker get repairs"
    $resp = Invoke-RestMethod "http://localhost:8080/api/repair/worker?worker_id=1" -Method Get
    Write-Host "Success: $($resp | ConvertTo-Json)"

    Write-Host "`n=== All tests passed ==="
} catch {
    Write-Host "Error: $($_.Exception.Message)"
}