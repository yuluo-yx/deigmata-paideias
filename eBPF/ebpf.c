int ebpf(void *ctx)
{
    bpf_trace_printk("hi, ebpf");
    return 0;
}
