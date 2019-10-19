# android-tv-remote-local

App for local network that controls Android TV by executing remote `adb` key event commands. Can be triggered by REST call from local network or by AWS SQS event.

##### Requires:
 * `adb` present in PATH
 * `~/.aws/credentials` specifying `aws_access_key_id` and `aws_secret_access_key`
 * Android TV IP address and port should be passed as first argument.
Example `./android-remote-tv "192.168.1.99:5555"`

##### Supported commands:
 * VOLUME_UP (N notches - optional)
 * VOLUME_DOWN (N notches - optional)
 * MUTE
 * POWER
 * SET_VOLUME N
