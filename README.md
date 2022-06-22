# Log Proxy Service

This is a sample application to view android crash log from device or from web browser, 
this is a web service, which can be run on the android device.

In order to run the service on the device, the debug mode need to be enabled, 
and compile the application and upload the application to the device, 
then use `adb shell` command to login to the device and start the service, 
after the service started, we can view the crash log from any web browser 
which was on the same network.

# Build application to android device

This is the build script on my machine.

## Emulator

```
GOARCH=amd64 GOOS=linux go build -o logproxyservice main.go
```

## Physical device

My Android is Redmi Note 9, the cpu info is `AArch64`, and the build script is:

```
GOARCH=arm64 GOOS=linux go build -o logproxyservice main.go
```

# Upload the application

```
adb -s <$ANDROID_SERIAL> push logproxyservice /data/local/tmp/
```

# Running the application

```
adb -s <$ANDROID_SERIAL> shell
cd /data/local/tmp

# Just run the application
./logproxyservice

# Run the application in background as a service
nohup ./logproxyservice > log_logproxyservice.log 2>&1 &

# The log will be saved on log_logproxservice.log, it will output the ip adress of the device
```

# View the crahlog from web browser

Open the web browser, and visit the service with the url `http://<ipaddress>:8080/crashlog`

If visit from the device, then just visit `http://localhost:8080`
