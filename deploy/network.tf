//network.tf

resource "aws_vpc" "network-go_ms" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "network-go_ms"
  }
}

resource "aws_eip" "network-go_ms" {
  instance = aws_instance.test-ec2-instance.id
  vpc      = true
}
