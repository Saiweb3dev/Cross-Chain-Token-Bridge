import FlickeringGrid from "@/components/magicui/flickering-grid";
import Button from "../ui/Button";
export function HeroSection() {
  return (
    <div className="relative h-screen rounded-lg w-full bg-background overflow-hidden border">
      <FlickeringGrid
        className="z-0 absolute inset-0 size-full w-full h-screen"
        squareSize={6}
        gridGap={6}
        color="#AD49E1"
        maxOpacity={0.2}
        flickerChance={0.2}
      />
      <div className="absolute inset-y-0 left-20 flex items-center z-10">
        <div className="max-w-4xl px-8 md:px-16 space-y-6">
          <h1 className="text-4xl md:text-6xl font-bold text-purple-700">
            Effortless Cross-Chain Token Transfers
          </h1>
          <span className="block text-lg md:text-xl text-gray-700">
            Seamlessly transfer your custom tokens between multiple blockchains with our secure, fast, and easy-to-use solution. Unlock the full potential of decentralized finance across chains.
          </span>
          <Button text="Get Started" href="/Contracts">
          </Button>
        </div>
      </div>
    </div>
  );
}