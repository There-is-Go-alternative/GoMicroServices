resource "aws_internet_gateway" "network-go_ms-gw" {
  vpc_id = aws_vpc.network-go_ms.id

  tags = {
    Name = "network-go_ms-gw"
  }

}