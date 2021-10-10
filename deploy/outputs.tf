output "instance_public_ip" {
#  value = module.ec2.*.public_ip
  value = aws_instance.test-ec2-instance.public_ip
}

output "instance_private_ip" {
  value = aws_instance.test-ec2-instance.private_ip
}

#output "instance_private_ip2" {
#  value = aws_vpc.network-go_ms-gw.
#}
