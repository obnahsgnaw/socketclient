interface Package {
    action: number,
    data: string
}

interface PackageCoder {
    Encode(pkg: Package): string;

    Decode(str: string): Package;
}