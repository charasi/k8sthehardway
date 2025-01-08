
def get_ip_address(filename):
    with open(filename, 'r') as f:
        first_line = f.readline().strip()
        return first_line
