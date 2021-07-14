# Firecrypt
Firecrypt is an application to encrypt, decrypt & manage Firefox profiles.

## Example Screenshot
![Example Screenshot](../assets/firecrypt_screenshot_1.png?raw=true)

## Usage
After executing Firecrypt, choose a profile from the list. If the profile has not already been configured, read the disclaimer and choose a password. Once configured, you can select the profile and enter your password to launch it. After you are finished using the profile, you can quit Firefox and Firecrypt will automatically reappear with the configuration options.

## Warning
This program is highly experimental! Cryptographic security is not a guarantee and data loss may occur.

## Important Notes
- Firecrypt currently only supports macOS.
- After opening a profile, Firecrypt will hide itself but will still be running. It will reappear once you quit Firefox.
- Firecrypt comes with _minimal_ safeguards in place. Do not attempt to break it as doing so will almost certainly work and cause data loss.
- To completely remove a profile from Firecrypt, you must unlock it and remove the `.__firecrypt_hash__` file in the profile directory. This file is hidden from the Finder so you should delete it using a terminal.
- If Firecrypt `panic`s during operation, you must check your Firefox profiles directory and remove any left over files (i.e. misplaced `.zip` files) before trying again.
- If you launch a Firefox profile while it is encrypted, you will receive the error "Your Firefox profile cannot be loaded. It may be missing or inaccessible."
- Firecrypt can currently only launch Firefox installations located at `/Applications/Firefox`.

## Disable Profile Already Open Check
When a profile is selected, Firecrypt will check to be sure it is not already being used. This process is slow and very inefficient. If you want to disable this safeguard for added speed, start Firecrypt with the `--no-check-profile-open` flag.

## To-Do
- [ ] Add support for other operating systems
- [ ] Add support for launching other Firefox installations
- [ ] Add support for profiles in other locations
- [ ] Add duress profiles

## License
[MIT](https://choosealicense.com/licenses/mit/)