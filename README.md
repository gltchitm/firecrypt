# Firecrypt
Firecrypt is an application to encrypt, decrypt & manage Firefox profiles.

## Screenshot
![Screenshot](../assets/firecrypt_screenshot_1.png?raw=true)

## Usage
After executing Firecrypt, choose a profile from the list. If the profile has not already been configured, read the disclaimer and choose a password. Once configured, you can select the profile and enter your password to launch it. After you are finished using the profile, you can quit Firefox and Firecrypt will automatically reappear with the configuration options.

## Warning
This program is highly experimental! Cryptographic security is not a guarantee and data loss may occur.

## Important Notes
- Firecrypt currently only supports macOS.
- After opening a profile, Firecrypt will hide itself but will still be running. It will reappear once you quit Firefox.
- To completely remove a profile from Firecrypt, you must unlock it and remove the `.__firecrypt_key__` file in the profile directory. This file is hidden from the Finder so you should delete it using a terminal.
- If you launch a Firefox profile while it is encrypted, you will receive the error "Your Firefox profile cannot be loaded. It may be missing or inaccessible."
- If you launch a Firefox profile while it is being configured in Firecrypt, you will receive the error "A copy of Firefox is already open. Only one copy of Firefox can be open at a time."
- Firecrypt can currently only launch Firefox installations located at `/Applications/Firefox`.

## Legacy Profiles
Legacy profiles (pre-V2) are no longer supported. You can use [this version of Firecrypt](https://github.com/gltchitm/firecrypt/tree/69b87376d4c7c8e05a6e0c4db1339d51b5a3002f) to migrate legacy profiles to Version 2, at which point they can be used in the latest version of Firecrypt.

## To-Do
- [ ] Add support for other operating systems
- [ ] Add support for launching other Firefox installations
- [ ] Add support for profiles in other locations
- [ ] Add duress profiles

## License
[MIT](LICENSE)
