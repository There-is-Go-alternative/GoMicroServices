variable "instance_name" {
  description = "Value of the Name tag for the EC2 instance"
  type        = string
  default     = "GoMicroService"
}

# https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/finding-an-ami.html#finding-quick-start-ami
variable "ami_name" {
  description = "Value of the Name for the AMI EC2 instance"
  type        = string
  #  default     = "Amazon Linux 2 AMI (HVM)"
  default = "Ubuntu Server 20.04 LTS (HVM)"
}

variable "ami_id" {
  description = "Value of the ID for the AMI EC2 instance"
  type        = string
  #  default     = "ami-072056ff9d3689e7b"
  default = "ami-0f7cd40eac2214b37"
}

variable "ami_key_pair_name" {
  description = "Name of the key pair the has the rights on amazon"
  type        = string
  default     = "deploy_key_thinkpad_popos"
}

variable "enable_monitoring" {
  description = "Name of the key pair the has the rights on amazon"
  type        = string
  default     = "deploy_key_thinkpad_popos"
}

