# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.
  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.define "router" do |router|
      router.vm.box = "ubuntu/focal64"
      router.vm.network :private_network, ip: "192.168.50.1", netmask:"255.255.255.0"
      router.vm.network :public_network, bridge: "en0: Wi-Fi (AirPort)"
      router.vm.provider :VirtualBox do |vb|
            vb.customize ["modifyvm", :id, "--nicpromisc1", "allow-all","--nicpromisc2", "allow-all"]
            end
      router.vm.provision "shell",
          run: "always",
          inline: <<-SHELL
          sysctl net.ipv4.ip_forward=1
          route add default gw 192.168.2.1 2>/dev/null || true
          sudo apt install -y python3-pip rustc openssl build-essential libssl-dev libffi-dev python3-dev python3
          pip3 install --upgrade setuptools  setuptools-rust
          pip3 install ansible
          cd go/src/covert
          ansible-playbook pb-router.yaml
#           iptables -t nat -A POSTROUTING -o enp0s9 -j MASQUERADE
#           iptables -A FORWARD -i enp0s9 -o enp0s8 -m state --state RELATED,ESTABLISHED -j ACCEPT
#           iptables -A FORWARD -i enp0s8 -o enp0s9 -j ACCEPT
         SHELL
  end
  config.vm.define "server" do |server|
      server.vm.box = "ubuntu/focal64"
      server.vm.network :private_network, ip: "192.168.50.2",netmask:"255.255.255.0"
      server.vm.provision "shell",
          run: "always",
          inline: <<-SHELL
          route add default gw 192.168.50.1 2>/dev/null || true
#           sudo apt update
#           sudo apt install -y python3-pip rustc openssl build-essential libssl-dev libffi-dev python3-dev python3
#           pip3 install --upgrade setuptools  setuptools-rust
#           pip3 install ansible
          SHELL
  end
  config.vm.define "client" do |client|
      client.vm.box = "ubuntu/focal64"
      client.vm.network :private_network, ip: "192.168.50.3",netmask:"255.255.255.0"
      client.vm.provision "shell",
           run: "always",
           inline: "route add default gw 192.168.50.1 2>/dev/null || true"
  end


  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # NOTE: This will enable public access to the opened port
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine and only allow access
  # via 127.0.0.1 to disable public access
  # config.vm.network "forwarded_port", guest: 80, host: 8080, host_ip: "127.0.0.1"

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
#   config.vm.network "private_network", ip: "192.168.50.100"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
#   config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  config.vm.synced_folder "./", "/home/vagrant/go/src/covert"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    vb.gui = false

    # Customize the amount of memory on the VM:
    vb.memory = "1024"
  end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Ansible, Chef, Docker, Puppet and Salt are also available. Please see the
  # documentation for more information about their specific syntax and use.

  config.vm.provision "shell", run: "always",inline: <<-SHELL
    sudo apt-get update -y
    sudo apt upgrade -y
    sudo apt install -y python3-pip rustc openssl build-essential libssl-dev libffi-dev python3-dev python3
    pip3 install --upgrade setuptools  setuptools-rust
    pip3 install ansible
    pip3 install ansible
    [ -d /home/vagrant/go/src ] || mkdir -p /home/vagrant/go/src
    chown -R vagrant:vagrant /home/vagrant/go
    cd /home/vagrant/go/src/covert && [ -f "go1.16.5.linux-amd64.tar.gz" ] || wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc
    apt-get install libpcap-dev -y
    export GOHOME=/home/vagrant/go/
    export GOROOT=/home/vagrant/go/
#     go build server.go
#     go build client.go
  SHELL

end
