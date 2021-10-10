terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }
  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region  = "eu-west-3"
}

data "aws_availability_zones" "current" {}

//servers.tf
resource "aws_instance" "test-ec2-instance" {
  ami           = var.ami_id
  instance_type = "t2.micro"
  # https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
  key_name = var.ami_key_pair_name
  monitoring = var.

  security_groups = [
    aws_security_group.ingress-all-test.id,
  ]

#  provisioner "remote-exec" {
#    inline = [
#      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
#      "sudo add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable'",
#      "sudo apt-get update",
#      "apt-cache policy docker-ce",
#      "sudo apt-get install -y docker-ce",
#      "sudo curl -L 'https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)' -o /usr/local/bin/docker-compose",
#      "sudo chmod +x /usr/local/bin/docker-compose",
#    ]
#    connection {
#      type = "ssh"
#      user = "ec2-user"
#      host = self.private_ip
#    }
#  }

#  provisioner "local-exec" {
#    inline = [
#      "curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
#      "sudo add-apt-repository 'deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable'",
#      "sudo apt-get update",
#      "apt-cache policy docker-ce",
#      "sudo apt-get install -y docker-ce",
#      "sudo curl -L 'https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)' -o /usr/local/bin/docker-compose",
#      "sudo chmod +x /usr/local/bin/docker-compose",
#    ]
#    connection {
#      type = "ssh"
#      user = "ec2-user"
#      host = self.private_ip
#    }
#  }

#  provisioner "file" {
#    source      = "../docker-compose.yml"
#    destination = "."
#
#    connection {
#      type = "ssh"
#      user = "ec2-user"
#      host = self.private_ip
#    }
#  }

  #  provisioner "remote-exec" {
  #    inline = [
  #      "touch hello.txt",
  #      "echo helloworld remote provisioner >> hello.txt",
  #    ]
  #  }

  user_data = file("install_docker.sh")

  tags = {
    Name        = var.instance_name
    Description = var.ami_name
  }

  subnet_id = aws_subnet.subnet-uno.id
}

#resource "aws_security_group" "instance" {
#  name = "go-ms-instance_sec"
#
#  ingress {
#    from_port   = 80
#    to_port     = 80
#    protocol    = "tcp"
#    cidr_blocks = ["0.0.0.0/0"]
#  }
#}
#
#resource "aws_launch_configuration" "test-ec2-instance" {
#  image_id      = var.ami_id
#  instance_type = "t2.micro"
#
#  security_groups = [aws_security_group.ingress-all-test.id]
#
#  user_data = file("install_docker.sh")
#
#  lifecycle {
#    create_before_destroy = true
#  }
#}

#resource "aws_launch_configuration" "app_server" {
#  image_id      = "ami-0f7cd40eac2214b37"
#  instance_type = "t2.micro"
#
#  security_groups = [aws_security_group.instance.id]
#
#  user_data = <<-EOF
#    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
#    sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
#    sudo apt-get update
#    apt-cache policy docker-ce
#    sudo apt-get install -y docker-ce
#    sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
#    sudo chmod +x /usr/local/bin/docker-compose
#    docker-compose build
#    docker-compose up -d
#    EOF
#
#  lifecycle {
#    create_before_destroy = true
#  }
#}
#
#resource "aws_instance" "app_server" {
#  ami           = "ami-0f7cd40eac2214b37"
#  instance_type = "t2.micro"
#  key_name= aws_key_pair.deployer.key_name
#
#  provisioner "remote-exec" {
#    inline = [
#      "touch hello.txt",
#      "echo helloworld remote provisioner >> hello.txt",
#    ]
#  }
#  connection {
#    type        = "ssh"
#    host        = self.public_ip
#    user        = "ubuntu"
#    private_key = file("/home/hobbit/.ssh/id_ed25519")
#    timeout     = "4m"
#  }
#
#  tags = {
#    Name = var.instance_name
#  }
#
#}
#
#resource "aws_key_pair" "deployer" {
#  key_name   = "id_ed25519"
#  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIYfE/HsdPgK2J3af9QNEJ3qDeuvSu3LYyGMdNq7/+Eu dbernard.dev@gmail.com"
#}
