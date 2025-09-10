import * as React from "react";
import {
  Star,
  Leaf,
  Zap,
  ShieldCheck,
  Sun,
  Globe2,
  Recycle,
  Quote,
} from "lucide-react";
import { Link } from "react-router";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ProductCard } from "@/components/ui/product-card";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

import { IProduct } from "@/types/product";

import { ProductService } from "@/services/product";

import solarFarmImg from "@/assets/solar-farm.webp";
import premiumSolarImg from "@/assets/premium-solar-panel.jpg";

function HomePage() {
  const [popProds, setPopProds] = React.useState<IProduct[]>([]);
  const [offerProds, setOfferProds] = React.useState<IProduct[]>([]);

  const handleAddToCart = (id: number) => {
    console.log("Added to cart:", id);
  };

  React.useEffect(() => {
    async function run() {
      const popProds = await ProductService.getProducts({
        avgscr: 4.0,
      });

      const offerProds = await ProductService.getProducts({
        offr: true,
        lim: 4,
        offst: 0,
      });

      setPopProds(() => popProds);
      setOfferProds(() => offerProds);
    }

    run();
  }, []);

  return (
    <main className="min-h-screen w-full">
      <section className="max-w-7xl mx-auto px-6 py-12 grid grid-cols-1 md:grid-cols-2 gap-8 items-center">
        <div className="space-y-6">
          <h1 className="text-4xl md:text-5xl font-extrabold leading-tight">
            Power Your Future with Clean Energy
          </h1>
          <p className="text-muted-foreground max-w-lg">
            We offer high-performance solar panels and sustainable energy
            solutions — designed to help you save money, reduce emissions, and
            create a cleaner planet.
          </p>
          <div className="flex gap-3 mt-4">
            <Button asChild>
              <Link to="/product">Shop Now!</Link>
            </Button>
            <Button variant="outline">Learn About Solar</Button>
          </div>

          <div className="mt-6 flex gap-6">
            <div className="flex items-center gap-2">
              <Star className="h-5 w-5 text-yellow-500" />
              <div>
                <div className="text-sm font-medium">4.9</div>
                <div className="text-xs text-muted-foreground">
                  Customer rating
                </div>
              </div>
            </div>

            <div className="flex items-center gap-2">
              <div className="text-sm font-medium">Free Consultations</div>
              <div className="text-xs text-muted-foreground">
                Personalized energy plan
              </div>
            </div>
          </div>
        </div>

        <div className="rounded-xl overflow-hidden shadow-lg">
          <img
            src={solarFarmImg}
            alt="Solar Panels"
            className="w-full h-full object-cover"
          />
        </div>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16 grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
        <div className="space-y-3">
          <Leaf className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">Eco-Friendly</h3>
          <p className="text-muted-foreground text-sm">
            Built with sustainable materials, our panels actively reduce your
            carbon footprint.
          </p>
        </div>

        <div className="space-y-3">
          <Zap className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">Efficient & Powerful</h3>
          <p className="text-muted-foreground text-sm">
            Generate maximum energy with minimal roof space using
            high-efficiency solar technology.
          </p>
        </div>

        <div className="space-y-3">
          <ShieldCheck className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">10+ Year Warranty</h3>
          <p className="text-muted-foreground text-sm">
            Weather-resistant panels with industry-leading guarantees for peace
            of mind.
          </p>
        </div>
      </section>

      <section className="max-w-7xl mx-auto px-6 pb-16 grid grid-cols-1 sm:grid-cols-3 gap-8 text-center">
        <div className="space-y-3">
          <Sun className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">Solar Made Simple</h3>
          <p className="text-muted-foreground text-sm">
            Easy installation process and clear guides make switching to solar
            effortless.
          </p>
        </div>

        <div className="space-y-3">
          <Globe2 className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">Global Impact</h3>
          <p className="text-muted-foreground text-sm">
            Every purchase contributes to global renewable energy initiatives.
          </p>
        </div>

        <div className="space-y-3">
          <Recycle className="mx-auto h-10 w-10 text-primary" />
          <h3 className="text-lg font-semibold">Recyclable Materials</h3>
          <p className="text-muted-foreground text-sm">
            Panels designed for end-of-life recycling to minimize waste.
          </p>
        </div>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold mb-6 text-center">
          Latest from the Blog
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <article className="rounded-lg border p-4 hover:shadow-md transition flex flex-col gap-3">
            <img
              src="https://picsum.photos/seed/prod-0/640/480"
              alt="Solar Trends 2025"
              className="rounded-lg object-cover"
            />
            <h3 className="font-semibold">Top Solar Trends for 2025</h3>
            <p className="text-muted-foreground text-sm">
              Discover where the solar industry is headed and how you can stay
              ahead.
            </p>
            <Button size="sm" variant="outline" className="mt-auto">
              Read More
            </Button>
          </article>

          <article className="rounded-lg border p-4 hover:shadow-md transition flex flex-col gap-3">
            <img
              src="https://picsum.photos/seed/prod-1/640/480"
              alt="Home Energy Guide"
              className="rounded-lg object-cover"
            />
            <h3 className="font-semibold">Ultimate Guide to Home Solar</h3>
            <p className="text-muted-foreground text-sm">
              Step-by-step advice on installing and optimizing your solar setup.
            </p>
            <Button size="sm" variant="outline" className="mt-auto">
              Read More
            </Button>
          </article>

          <article className="rounded-lg border p-4 hover:shadow-md transition flex flex-col gap-3">
            <img
              src="https://picsum.photos/seed/prod-2/640/480"
              alt="Battery Storage Tips"
              className="rounded-lg object-cover"
            />
            <h3 className="font-semibold">Battery Storage Made Easy</h3>
            <p className="text-muted-foreground text-sm">
              Learn how to maximize energy storage and power your home
              overnight.
            </p>
            <Button size="sm" variant="outline" className="mt-auto">
              Read More
            </Button>
          </article>
        </div>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold mb-8 text-center">
          Popular Solar Products
        </h2>
        <Carousel className="w-full py-6 px-8 max-[440px]:px-0">
          <CarouselContent className="items-stretch gap-5">
            {popProds.map((p) => (
              <CarouselItem className="min-[728px]:basis-1/2 min-[1090px]:basis-1/3 min-[1330px]:basis-1/4">
                <ProductCard
                  productId={p.id}
                  image={
                    p.mainImage
                      ? `http://localhost:5000/product/image/${p.mainImage.imageName}`
                      : ""
                  }
                  storeName={p.store.name}
                  title={p.name}
                  categoryName={p.subcategory.name}
                  price={p.price}
                  score={p.averageScore}
                  onAddToCart={handleAddToCart}
                  isLowStock={p.totalQuantity < 5}
                  isOutStock={p.totalQuantity === 0}
                  discount={p.offer != null ? p.offer.discount : 0}
                />
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious className="-left-3 max-[440px]:hidden" />
          <CarouselNext className="-right-3 max-[440px]:hidden" />
        </Carousel>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold mb-8 text-center">
          How Going Solar Works
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
          <div className="space-y-3">
            <div className="text-4xl font-bold text-primary">1</div>
            <h3 className="text-lg font-semibold">Free Consultation</h3>
            <p className="text-muted-foreground text-sm">
              We evaluate your home and energy needs to design the perfect solar
              solution.
            </p>
          </div>

          <div className="space-y-3">
            <div className="text-4xl font-bold text-primary">2</div>
            <h3 className="text-lg font-semibold">Professional Installation</h3>
            <p className="text-muted-foreground text-sm">
              Our certified team handles everything — from permits to final
              setup.
            </p>
          </div>

          <div className="space-y-3">
            <div className="text-4xl font-bold text-primary">3</div>
            <h3 className="text-lg font-semibold">Start Saving</h3>
            <p className="text-muted-foreground text-sm">
              Begin generating clean power and watch your monthly bills shrink.
            </p>
          </div>
        </div>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold mb-8 text-center">
          Featured Solar Solutions
        </h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="rounded-xl overflow-hidden shadow-lg border flex flex-col">
            <img
              src={premiumSolarImg}
              alt="Premium Solar Panel"
              className="object-cover w-full aspect-video"
            />
            <div className="p-6 flex flex-col gap-2 h-full">
              <h3 className="text-xl font-semibold">Premium Solar Panel X</h3>
              <p className="text-muted-foreground text-sm">
                Our best-selling high-performance panel, perfect for residential
                rooftops. Exercitationem officia quia corporis doloribus
                voluptatum. Dolor dolorem nesciunt fuga dignissimos. Nam
                accusantium praesentium autem aut odio sint repudiandae.
                Incidunt veritatis eos quam doloremque tempore cupiditate
                cumque.
              </p>
              <Button onClick={() => handleAddToCart(1)} className="mt-auto">
                Add to Cart
              </Button>
            </div>
          </div>

          <div className="w-full @container">
            <div className="grid grid-cols-1 gap-6 @md:grid-cols-2">
              {offerProds.map((p) => (
                <ProductCard
                  key={p.id}
                  productId={p.id}
                  image={
                    p.mainImage
                      ? `http://localhost:5000/product/image/${p.mainImage.imageName}`
                      : ""
                  }
                  storeName={p.store.name}
                  title={p.name}
                  categoryName={p.subcategory.name}
                  price={p.price}
                  score={p.averageScore}
                  onAddToCart={handleAddToCart}
                  isLowStock={p.totalQuantity < 5}
                  isOutStock={p.totalQuantity === 0}
                  discount={p.offer != null ? p.offer.discount : 0}
                />
              ))}
            </div>
          </div>
        </div>
      </section>

      <Separator />

      <section className="max-w-7xl mx-auto px-6 py-16">
        <h2 className="text-3xl font-bold text-center mb-10">
          What Our Customers Say
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[
            {
              name: "Sarah M.",
              quote:
                "Switching to solar with this company was the best decision we've made. Our energy bills are half of what they used to be!",
            },
            {
              name: "James L.",
              quote:
                "The installation was quick and professional. The panels look great on our roof and perform even better.",
            },
            {
              name: "Priya S.",
              quote:
                "Customer service is amazing! They walked me through every step and made the process stress-free.",
            },
          ].map((t, i) => (
            <div
              key={i}
              className="rounded-xl border p-6 text-center shadow-sm"
            >
              <Quote className="mx-auto text-primary mb-4 h-6 w-6" />
              <p className="text-sm text-muted-foreground italic mb-4">
                “{t.quote}”
              </p>
              <p className="font-semibold">{t.name}</p>
            </div>
          ))}
        </div>
      </section>

      <Separator />

      <section className="bg-primary text-primary-foreground py-16 text-center">
        <h2 className="text-3xl font-bold mb-4">Ready to Go Solar?</h2>
        <p className="max-w-2xl mx-auto mb-6">
          Take control of your energy future today. Request a free consultation
          and see how much you can save with solar power.
        </p>
        <Button size="lg" variant="secondary">
          Get Your Free Quote
        </Button>
      </section>
    </main>
  );
}

export default HomePage;
