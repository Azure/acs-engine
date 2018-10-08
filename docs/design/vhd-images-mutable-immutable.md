# Problem statement

There is a constraint on the number of versions per VHD VM SKU.

# Current usage pattern of VHD VM SKU publishing

We increment and publish a new versioned image if any of the following is true:

1. A new dependent artifact is introduced into the provisioning flow (e.g., a new Kubernetes version --> new hyperkube image)
2. A new OS layer configuration implementation (e.g., Azure CNI binary moves from one location to another on disk)
3. An improved/optimized implementation of artifact delivery is shipped (artifact delivery implementation is shared between pre-baked [VHD] and runtime [CSE] provision flows to facilitate operational resilience and accommodate hotfix velocity)
  - In other words, if we need to release a zero day security patch we are able to to that via CSE without having to wait for VHD publication.

# Advantages of immutable image references (e.g., `0.18.0`) in provisioning implementation

- Inherently tightly couples the VHD implementation to the OS image implementation
  - This is an advantage insofar as the VHD imlementation is curated specifically for AKS. This allows a predictable, hygienic flow thusly:
    - customer request -->
    - --> RP -->
    - --> acs-engine -->
    - --> VHD -->
    - --> VM provisioning runtime
  - In the above flow, the RP + acs-engine + VHD + VM provisioning runtime will all have been thoroughly tested @ distinct versions _as an operational unit_ before being exposed to AKS customers. In other words, in terms of provisioning SLA, the operation version being released/exposed to customers is a distinct composed version that can be expressed as *RP+acs-engine+VHD*
- Is what we're doing at present (carried over from base Ubuntu 16.04 LTS implementation); requires no additional re-architecting

# Disadvantages of immutable image references

- Because there is re-use between the VHD and non-VHD CSE implementation, we often absorb the overhead of having to publish new VHD images when we optimize the CSE delivery/installation surface area.
- More storage overhead due to more images

# Advantages of mutable image references (e.g., `latest`) in provisioning implementation

- Decomposes the VHD implementation from acs-engine runtime implementation, allowing for re-use patterns
  - There is no use-case at present, but we can hypothesize:
    - General purpose "hyperkube-enabled" OS images that non-acs-engine Kubernetes-interested projects can opt into
- Less storage overhead due to fewer images

# Disadvantages of mutable image references

- Requires an architectural refresh, as the current VHD implementation is a forward iteration of the prior Ubuntu 16.04 LTS immutable image usage pattern
- Additional CI overhead, as we will have to "promote" images from a "staging" mutable image reference to a production "latest" image
- Introduces additional troubleshooting overhead when diagnosing production issues in the provisioning stack
  - e.g., "Which "latest" was this vm provisioned with???"