#pragma semicolon 1
#pragma newdecls required
#include <sourcemod>

public Plugin myinfo = { name = "D2C Frame Metrics", author = "DotaClassic", description = "Read-only rolling wall-clock frame interval telemetry", version = "1.0.0", url = "" };
#define CAPACITY 4096
float g_Time[CAPACITY];
float g_Interval[CAPACITY];
float g_Last;
int g_Next;
int g_Count;

public void OnPluginStart() { RegServerCmd("d2c_frame_stats", FrameStats); }
public void OnMapStart() { g_Last = 0.0; g_Next = 0; g_Count = 0; }
public void OnGameFrame() {
    float now = GetEngineTime();
    if (g_Last > 0.0 && now > g_Last) {
        g_Time[g_Next] = now;
        g_Interval[g_Next] = (now - g_Last) * 1000.0;
        g_Next = (g_Next + 1) % CAPACITY;
        if (g_Count < CAPACITY) g_Count++;
    }
    g_Last = now;
}

public Action FrameStats(int args) {
    float now = GetEngineTime();
    float sum = 0.0;
    float budget = GetTickInterval() * 1500.0;
    int over = 0;
    ArrayList samples = new ArrayList();
    for (int i = 0; i < g_Count; i++) {
        if (now - g_Time[i] <= 5.0) {
            float value = g_Interval[i];
            samples.Push(view_as<int>(value));
            sum += value;
            if (value > budget) over++;
        }
    }
    int n = samples.Length;
    if (n == 0) { PrintToServer("D2C_FRAME frames=0 mean_ms=0 p95_ms=0 p99_ms=0 max_ms=0 over_ratio=0"); delete samples; return Plugin_Handled; }
    SortADTArray(samples, Sort_Ascending, Sort_Float);
    float p95 = view_as<float>(samples.Get(RoundToCeil(float(n) * 0.95) - 1));
    float p99 = view_as<float>(samples.Get(RoundToCeil(float(n) * 0.99) - 1));
    float max = view_as<float>(samples.Get(n - 1));
    PrintToServer("D2C_FRAME frames=%d mean_ms=%.3f p95_ms=%.3f p99_ms=%.3f max_ms=%.3f over_ratio=%.6f", n, sum / float(n), p95, p99, max, float(over) / float(n));
    delete samples;
    return Plugin_Handled;
}
