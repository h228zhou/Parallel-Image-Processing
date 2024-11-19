import subprocess
import matplotlib.pyplot as plt

def run_benchmark(patterns, data_sets, thread_counts):
    """
    Runs the benchmark for each combination of pattern, data set, and thread count.
    Returns dictionaries containing the average times for sequential and parallel runs.
    """
    seq_times = {}
    parallel_times = {pattern: {data_set: {} for data_set in data_sets} for pattern in patterns if pattern != "sequential"}

    for pattern in patterns:
        print(f"Running pattern: {pattern}")
        for data_set in data_sets:
            if pattern == "sequential":
                cmd = f"go run ../editor/editor.go {data_set}"
                try:
                    output = subprocess.check_output(cmd, shell=True).strip()
                    time = float(output)
                    seq_times[data_set] = time
                    print(f"Sequential run - Data set: {data_set}, Time: {time:.4f}s")
                except subprocess.CalledProcessError as e:
                    print(f"Error running command: {cmd}")
                    print(e.output.decode())
            else:
                for thread_count in thread_counts:
                    cmd = f"go run ../editor/editor.go {data_set} {pattern} {thread_count}"
                    try:
                        output = subprocess.check_output(cmd, shell=True).strip()
                        time = float(output)
                        parallel_times[pattern][data_set][thread_count] = time
                        print(f"Parallel run - Pattern: {pattern}, Data set: {data_set}, Threads: {thread_count}, Time: {time:.4f}s")
                    except subprocess.CalledProcessError as e:
                        print(f"Error running command: {cmd}")
                        print(e.output.decode())

    return seq_times, parallel_times

def compute_speedups(seq_times, parallel_times, thread_counts):
    """
    Computes speedups for each pattern, data set, and thread count.
    Returns a nested dictionary containing speedup values.
    """
    speedups = {pattern: {data_set: [] for data_set in seq_times} for pattern in parallel_times}

    for pattern, data_sets in parallel_times.items():
        for data_set, times in data_sets.items():
            seq_time = seq_times[data_set]
            for thread_count in thread_counts:
                par_time = times.get(thread_count)
                if par_time:
                    speedup = seq_time / par_time
                    speedups[pattern][data_set].append((thread_count, speedup))
                else:
                    print(f"No data for Pattern: {pattern}, Data set: {data_set}, Threads: {thread_count}")

    return speedups

def plot_speedups(speedups, thread_counts):
    """
    Plots speedup vs. number of threads for each data set and pattern.
    """
    for pattern, data_sets in speedups.items():
        plt.figure(figsize=(10, 8))
        for data_set, values in data_sets.items():
            threads = [t for t, _ in values]
            speedup_values = [s for _, s in values]
            plt.plot(threads, speedup_values, marker='o', linestyle='-', label=data_set.capitalize())

        plt.title(f'Speedup vs. Number of Threads for Pattern: {pattern.capitalize()}')
        plt.xlabel('Number of Threads')
        plt.ylabel('Speedup')
        plt.xticks(thread_counts)
        max_speedup = max([s for data in data_sets.values() for _, s in data])
        plt.yticks([i * 0.5 for i in range(int(2 * max_speedup) + 2)])
        plt.ylim(0, max_speedup + 1)
        plt.legend()
        plt.grid(True)
        plt.savefig(f'speedup_{pattern}.png')
        plt.show()

def main():
    # Define parameters
    patterns = ["s", "parslices", "parfiles"]
    data_sets = ["small", "mixture", "big"]
    thread_counts = [2, 4, 6, 8, 12]

    # Run benchmarks
    seq_times, parallel_times = run_benchmark(patterns, data_sets, thread_counts)

    # Compute speedups
    speedup_results = compute_speedups(seq_times, parallel_times, thread_counts)

    # Plot results
    plot_speedups(speedup_results, thread_counts)

if __name__ == '__main__':
    main()
