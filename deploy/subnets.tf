//subnets.tf

resource "aws_subnet" "subnet-uno" {
  cidr_block        = cidrsubnet(aws_vpc.network-go_ms.cidr_block, 3, 1)
  vpc_id            = aws_vpc.network-go_ms.id
  availability_zone = "eu-west-3a"
}

//subnets.tf
resource "aws_route_table" "route-table-test-env" {
  vpc_id = aws_vpc.network-go_ms.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.network-go_ms-gw.id
  }
  tags = {
    Name = "test-env-route-table"
  }
}

resource "aws_route_table_association" "subnet-association" {
  subnet_id      = aws_subnet.subnet-uno.id
  route_table_id = aws_route_table.route-table-test-env.id
}