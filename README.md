# Kai Suite

```What is the purpose Kai Suite ?
A pc suite for KaiOS device to manage events(Google Calendar/TODO) & contacts(locals/Google People)
and capability to send or received sms via desktop app.
```

## Kai Suite KaiOS client
https://github.com/arma7x/kai-suite-client

## Status[WIP]

### Guides(Disclaimer: Please backup your contacts before testing)

#### Connection
```
- Use ifconfig(linux) or ipconfig(windows) to get your wi-fi ip address
- Please setup port forwarding, if your pc/laptop not connected to KaiOS hotspot
```

#### Local Contacts
```
- The origin of contact is KaiOS Device/VCF
- Please use Restore, if you accidentally delete any contacts on your device
or when the KaiOS device is connected to Kai Suite for the first time
```

#### Google Contacts
```
- The origin of contact is Google People API
- Please use Restore, if you accidentally delete any contacts on yourdevice
or when the KaiOS device is connected to Kai Suite for the first time
```

#### Setup Google API(https://youtu.be/Wk6pk-uRUOE)
```
1. Create new project, visit https://console.cloud.google.com/
2. Enable People API & Calendar API
3. Configure Consent Screen
4. Create Credentials
5. Download the credential json file and rename it as credentials.json
6. Open credentials.json, search for http://localhost and replace it with urn:ietf:wg:oauth:2.0:oob
7. The credentials.json & Kai Suite(binary file) must reside in same folder/directory
```

#### Written in Go, powered by [Fyne](https://github.com/fyne-io/fyne)
