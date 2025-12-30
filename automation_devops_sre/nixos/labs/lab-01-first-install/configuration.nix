# Lab 1: First NixOS Configuration
# 
# This is your system configuration. It declares what your system should look like.
# Every line describes something about your operating system.
# When you run `nixos-rebuild switch`, NixOS reads this file and makes it real.
#
# Key principle: This file is the SOURCE OF TRUTH for your system.
# If you want to know what's installed, what services are running,
# or how things are configured - it's all here.

{ config, pkgs, ... }:

{
  ############################################################################
  # IMPORTS
  # 
  # Most NixOS systems need hardware-specific configuration
  # (bootloader type, disk layout, CPU features, etc.)
  # 
  # The nixos-generate-config tool creates hardware-configuration.nix
  # You should almost NEVER edit this manually - it's auto-generated
  # Just import it here
  ############################################################################
  
  imports =
    [
      ./hardware-configuration.nix  # Auto-generated, describes your hardware
    ];


  ############################################################################
  # BOOTLOADER CONFIGURATION
  #
  # This section tells NixOS how to boot your system.
  # It sets up GRUB bootloader and where to install it.
  #
  # Key concept: NixOS manages bootloader config automatically
  # (unlike traditional Linux where you edit grub.cfg manually)
  ############################################################################
  
  boot.loader.grub.enable = true;
  # Use GRUB 2 bootloader (most common)
  
  boot.loader.grub.device = "/dev/sda";
  # Install GRUB to this disk (MBR mode)
  # For EFI systems, use: boot.loader.systemd-boot.enable = true;
  
  boot.loader.grub.copyKernels = true;
  # Copy kernels to /boot partition (ensures redundancy)


  ############################################################################
  # KERNEL CONFIGURATION
  #
  # You can choose different kernels and set boot parameters
  ############################################################################
  
  boot.kernelPackages = pkgs.linuxPackages_latest;
  # Use latest stable kernel (other options: linuxPackages_6_1, etc.)


  ############################################################################
  # SYSTEM IDENTIFICATION
  #
  # These are basic identifiers for your system
  ############################################################################
  
  networking.hostName = "nixos-lab-01";
  # Your system's hostname (what appears in terminal prompt)
  # Change this to something meaningful for your system
  
  networking.domain = "local";
  # Domain name (typically "local" for home networks)


  ############################################################################
  # NETWORKING
  #
  # Configure network interfaces, DNS, firewall
  ############################################################################
  
  networking.useDHCP = lib.mkDefault true;
  # Use DHCP to automatically get IP address
  # Alternative: staticIp with manual IP configuration
  
  # For Wi-Fi:
  # networking.wireless.enable = true;
  # networking.wireless.networks = {
  #   "MyWiFiSSID".psk = "password";
  # };
  
  networking.firewall.enable = true;
  # Enable firewall (blocks everything by default, you whitelist)
  # Most home systems can leave this enabled
  
  # networking.firewall.allowedTCPPorts = [ 22 80 443 ];
  # Uncomment to open specific ports
  # (22=SSH, 80=HTTP, 443=HTTPS)


  ############################################################################
  # TIME ZONE
  #
  # Set your system timezone
  ############################################################################
  
  time.timeZone = "UTC";
  # Change to your timezone: "America/New_York", "Europe/London", etc.
  # Full list: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones


  ############################################################################
  # SYSTEM PACKAGES
  #
  # These packages are available to ALL users on the system.
  # Think of these as system-wide utilities everyone might need.
  #
  # Key point: Installing a package here doesn't pollute your system
  # Each package is isolated in /nix/store and linked cleanly
  ############################################################################
  
  environment.systemPackages = with pkgs; [
    # Essential tools (most systems have these)
    vim                 # Text editor
    git                 # Version control
    curl                # HTTP client
    wget                # Download utility
    htop                # System monitor (like top but nicer)
    
    # Utilities
    tree                # Directory structure viewer
    jq                  # JSON query tool
    fzf                 # Fuzzy finder
    
    # Add your own packages here!
    # Find packages at: https://search.nixos.org/packages
    # Examples: 
    # - nodejs (JavaScript runtime)
    # - python3 (Python programming)
    # - docker (container runtime)
  ];
  
  # If you want a package but don't know the exact name:
  # $ nix search nixpkgs "partial-name"
  # Example: $ nix search nixpkgs "node"


  ############################################################################
  # ENVIRONMENT VARIABLES
  #
  # Set system-wide environment variables
  # Available to all users in all shells
  ############################################################################
  
  environment.variables = {
    EDITOR = "vim";
    # When programs ask for your editor, use vim
    
    # Add your own:
    # CUSTOM_VAR = "value";
  };


  ############################################################################
  # SHELL CONFIGURATION
  #
  # NixOS can manage your default shell
  # Supports bash, zsh, fish, etc.
  ############################################################################
  
  programs.bash.enable = true;
  # Bash is usually enabled by default
  
  # For zsh:
  # programs.zsh.enable = true;
  # environment.shells = with pkgs; [ zsh ];


  ############################################################################
  # SERVICES
  #
  # Services are long-running processes (daemons) that start at boot.
  # NixOS manages their systemd units automatically.
  #
  # Key concept: Just enable them here, NixOS handles everything else
  # (config file generation, permissions, startup, etc.)
  ############################################################################
  
  # Example: SSH Server (useful for remote access)
  services.openssh.enable = false;
  # Uncomment to enable SSH server
  # Then: services.openssh.listenAddresses = [ "0.0.0.0" ];
  # And configure firewall to allow port 22
  
  # Example: NTP (network time sync)
  services.chrony.enable = true;
  # Keep system time synchronized with NTP servers


  ############################################################################
  # USERS
  #
  # Declare user accounts here.
  # NixOS will create/update them to match this configuration.
  #
  # Key principle: Users are part of your system declaration
  # Reproducible, version-controlled, atomic
  ############################################################################
  
  # Disable the root user login password (already disabled by default)
  users.mutableUsers = false;
  # When false: users can only be created/modified via this config
  # When true: you can also modify with useradd/passwd (more traditional)
  
  # For learning, let's keep it simple
  users.mutableUsers = true;
  # This lets you use `passwd` to set passwords
  
  # Create specific users:
  # users.users.alice = {
  #   isNormalUser = true;
  #   home = "/home/alice";
  #   shell = pkgs.bash;
  #   groups = [ "wheel" ];  # wheel group = sudo access
  #   initialPassword = "changeme";
  # };


  ############################################################################
  # SUDO CONFIGURATION
  #
  # Configure sudo (privilege escalation)
  ############################################################################
  
  security.sudo.enable = true;
  # Enable sudo (required for system admin tasks)


  ############################################################################
  # LOCALE AND INTERNACIONALIZATION
  #
  # Set language, keyboard layout, etc.
  ############################################################################
  
  i18n.defaultLocale = "en_US.UTF-8";
  # Set system language to US English, UTF-8 encoding
  
  console = {
    # Keyboard layout for login console (not GUI)
    keyMap = "us";
    # Change to "uk", "dvorak", etc. as needed
  };


  ############################################################################
  # SYSTEM STATE VERSION
  #
  # This is crucial! It tells NixOS which version of config format to use.
  # NEVER change this unless migrating systems (can break things).
  #
  # When you first install, nixos-generate-config sets this for you.
  # You should keep it as is for that installation.
  ############################################################################
  
  system.stateVersion = "23.11";
  # This says "this system was created with NixOS 23.11"
  # NixOS uses this for compatibility decisions
  # SEE RELEASE NOTES BEFORE CHANGING THIS
  
  # When to change stateVersion:
  # - Never for normal operation (keep it as generated)
  # - Only if migrating a system and reading docs carefully
  # - Changing incorrectly can break automatic migrations


  ############################################################################
  # OPTIONAL: DOCUMENTATION
  #
  # Configure what documentation is included
  ############################################################################
  
  documentation.enable = true;
  # Include man pages and documentation
  
  documentation.man.enable = true;
  # Enable man pages ($ man <command>)


  ############################################################################
  # OPTIONAL: SOUND/AUDIO
  #
  # For desktop systems with audio
  ############################################################################
  
  # sound.enable = true;  # Uncomment for sound support
  # hardware.pulseaudio.enable = true;  # Or use PulseAudio
  # OR use PipeWire (modern):
  # services.pipewire.enable = true;


  ############################################################################
  # END OF CONFIGURATION
  #
  # The structure above is typical for most NixOS systems.
  # You can:
  # - Enable more services (see search.nixos.org/options)
  # - Add more packages
  # - Create additional users
  # - Configure specific programs
  #
  # Each change = edit this file = $ nixos-rebuild switch = running system
  ############################################################################

}

# HELPFUL COMMANDS:
# 
# $ sudo nixos-rebuild switch
#   Apply configuration changes (rebuilds what changed)
#   
# $ sudo nixos-rebuild dry-build
#   Preview changes without applying them
#   
# $ sudo nixos-rebuild list-generations
#   See all system versions (rollback options)
#   
# $ sudo nixos-rebuild switch --rollback
#   Revert to previous system configuration
#   
# $ nix search nixpkgs "search-term"
#   Find package by name
#   
# $ nix-store -q --tree /run/current-system
#   See dependency tree of your entire system
#
# SEARCHING FOR OPTIONS:
# 
# Visit: https://search.nixos.org/options
# Search for configuration options (e.g., "ssh", "firewall", etc.)
# Each option shows:
# - What it does
# - Valid values
# - Default value
# - Example configurations
#
# KEY INSIGHT:
#
# Every line in this file has consequences but is completely safe.
# You can:
# - Add packages without breaking anything
# - Remove packages instantly
# - Change settings and rollback if needed
# - Test changes before applying to production
#
# This is how NixOS provides safety and reproducibility.
