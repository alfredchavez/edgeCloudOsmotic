use structopt::StructOpt;

#[derive(Debug, StructOpt)]
struct Opt {
    /// Evaluate HQ9+ source code
    #[structopt(short = "e", long = "eval")]
    source_code: Option<String>,
}

fn simple_sieve(limit: usize) -> usize {
    let mut is_prime = vec![true; limit + 1];
    is_prime[0] = false;
    if limit >= 1 { is_prime[1] = false }

    for num in 2..limit + 1 {
        if is_prime[num] {
            let mut multiple = num * num;
            while multiple <= limit {
                is_prime[multiple] = false;
                multiple += num;
            }
        }
    }

    is_prime.iter().enumerate()
        .filter_map(|(pr, &is_pr)| if is_pr { Some(pr) } else { None })
        .count()
}

fn main() {
    let opt = Opt::from_args();
    if let Some(src) = opt.source_code {
        let num = src.parse::<usize>().unwrap();
        println!("{:?}", simple_sieve(num));
    }

}