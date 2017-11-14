function Get-VmComputeNativeMethods()
{
        $signature = @'
                     [DllImport("vmcompute.dll")]
                     public static extern void HNSCall([MarshalAs(UnmanagedType.LPWStr)] string method, [MarshalAs(UnmanagedType.LPWStr)] string path, [MarshalAs(UnmanagedType.LPWStr)] string request, [MarshalAs(UnmanagedType.LPWStr)] out string response);
'@

    # Compile into runtime type
    Add-Type -MemberDefinition $signature -Namespace VmCompute.PrivatePInvoke -Name NativeMethods -PassThru
}

function New-HnsNetwork
{
    param
    (
        [parameter(Mandatory=$false, Position=0)]
        [string] $JsonString,
        [ValidateSet('L2Bridge')]
        [parameter(Mandatory = $false, Position = 0)]
        [string] $Type,
        [parameter(Mandatory = $false)] [string] $Name,
        [parameter(Mandatory = $false)] [string] $AddressPrefix,
        [parameter(Mandatory = $false)] [string] $Gateway,
        [parameter(Mandatory = $false)] [string] $DNSServer
    )

    Begin {
        if (!$JsonString) {
            $netobj = @{
                Type          = $Type;
            };

            if ($Name) {
                $netobj += @{
                    Name = $Name;
                }
            }

            if ($AddressPrefix -and  $Gateway) {
                $netobj += @{
                    Subnets = @(
                        @{
                            AddressPrefix  = $AddressPrefix;
                            GatewayAddress = $Gateway;
                        }
                    );
                }
            }

            if ($DNSServerName) {
                $netobj += @{
                    DNSServerList = $DNSServer
                }
            }

            $JsonString = ConvertTo-Json $netobj -Depth 10
        }

    }
    Process {
        $Method = "POST";
        $hnsPath = "/networks/";
        $request = $JsonString;

        Write-Verbose "Invoke-HNSRequest Method[$Method] Path[$hnsPath] Data[$request]"

        $output = "";
        $response = "";
        $hnsApi = Get-VmComputeNativeMethods
        $hnsApi::HNSCall($Method, $hnsPath, "$request", [ref] $response);

        Write-Verbose "Result: $response"
        if ($response)
        {
            try {
                $output = ($response | ConvertFrom-Json);
            } catch {
                Write-Error $_.Exception.Message
                return ""
            }
            if ($output.Error)
            {
                 Write-Error $output;
            }
            $output = $output.Output;
        }

        return $output;
    }
}

Export-ModuleMember -Function New-HNSNetwork
