import {
  Badge,
  Button,
  Card,
  Flex,
  Heading,
  Skeleton,
  Text,
} from "@radix-ui/themes";
import { useEffect, useState } from "react";
import {
  getJSON,
  type IndexerStatus,
  type Market,
  type PositionsPage,
  type PositionSummary,
} from "./api";
import {
  formatSyncTime,
  formatE18,
  groupInteger,
  percentFromRatio,
  shortAddress,
} from "./format";

interface Resource<T> {
  data: T | null;
  error: Error | null;
}

const POSITION_PAGE_SIZE = 10;
const POSITION_EMOJIS = ["🦊", "🐼", "🐸", "🐙", "🦄", "🐝", "🦉", "🐳", "🦜", "🦔"];

function useResource<T>(path: string): Resource<T> {
  const [resource, setResource] = useState<Resource<T>>({
    data: null,
    error: null,
  });

  useEffect(() => {
    let active = true;
    let controller: AbortController | null = null;

    const load = async () => {
      controller?.abort();
      controller = new AbortController();
      try {
        const data = await getJSON<T>(path, controller.signal);
        if (active) setResource({ data, error: null });
      } catch (error) {
        if (active && !(error instanceof DOMException && error.name === "AbortError")) {
          setResource((current) => ({
            data: current.data,
            error: error instanceof Error ? error : new Error("Unknown API error"),
          }));
        }
      }
    };

    void load();
    const timer = window.setInterval(load, 2_000);
    return () => {
      active = false;
      controller?.abort();
      window.clearInterval(timer);
    };
  }, [path]);

  return resource;
}

function App() {
  const status = useResource<IndexerStatus>("/api/v1/status");
  const market = useResource<Market>("/api/v1/market");
  const [cursors, setCursors] = useState<Array<string | null>>([null]);
  const cursor = cursors[cursors.length - 1];
  const positions = useResource<PositionsPage>(
    `/api/v1/positions?limit=${POSITION_PAGE_SIZE}${cursor ? `&cursor=${encodeURIComponent(cursor)}` : ""}`,
  );

  return (
    <main>
      <div className="container page-body">
        <Heading as="h1" size="7" mb="5">
          Market
        </Heading>
        <StatsCard market={market.data} status={status.data} positions={positions.data} />
        <PositionsFeed
          resource={positions}
          page={cursors.length}
          onPrevious={() => setCursors((current) => current.slice(0, -1))}
          onNext={() => {
            const nextCursor = positions.data?.nextCursor;
            if (nextCursor) setCursors((current) => [...current, nextCursor]);
          }}
        />
      </div>
    </main>
  );
}

function StatsCard({
  market,
  status,
  positions,
}: {
  market: Market | null;
  status: IndexerStatus | null;
  positions: PositionsPage | null;
}) {
  const liquidity = market
    ? (BigInt(market.totalSupplyAssets) - BigInt(market.totalBorrowAssets)).toString()
    : null;
  return (
    <Card className="stats-card">
      <Stat
        icon="S"
        label="Total supply"
        value={market ? formatE18(market.totalSupplyAssets) : null}
      />
      <Stat
        icon="B"
        label="Total borrow"
        value={market ? formatE18(market.totalBorrowAssets) : null}
      />
      <Stat icon="P" label="Positions" value={positions ? groupInteger(positions.total) : null} />
      <Stat icon="L" label="Available liquidity" value={liquidity ? formatE18(liquidity) : null} />
      <Stat
        icon="U"
        label="Utilization"
        value={market ? percentFromRatio(market.totalBorrowAssets, market.totalSupplyAssets) : null}
      />
      <Stat
        icon="#"
        label="Latest block"
        value={status?.chainHead ? groupInteger(status.chainHead) : null}
        detail={formatSyncTime(status?.lastSuccessfulSync ?? null)}
      />
    </Card>
  );
}

function Stat({
  icon,
  label,
  value,
  detail,
}: {
  icon: string;
  label: string;
  value: string | null;
  detail?: string;
}) {
  return (
    <div className="stat">
      <span className="stat-icon">{icon}</span>
      <div>
        <Text as="p" size="1" color="gray">{label}</Text>
        {value ? (
          <Text as="p" size="4" weight="bold" className="stat-value">
            {value}
          </Text>
        ) : (
          <Skeleton width="90px" height="24px" />
        )}
        {detail && (
          <Text as="p" size="1" color="gray">
            {detail}
          </Text>
        )}
      </div>
    </div>
  );
}

function PositionsFeed({
  resource,
  page,
  onPrevious,
  onNext,
}: {
  resource: Resource<PositionsPage>;
  page: number;
  onPrevious: () => void;
  onNext: () => void;
}) {
  return (
    <Card id="positions" className="feed-card">
      <Flex align="center" justify="between" className="feed-header">
        <Heading as="h2" size="4">Positions</Heading>
        {resource.data && (
          <Badge color="gray" variant="outline">
            Block #{groupInteger(resource.data.checkpoint.blockNumber)}
          </Badge>
        )}
      </Flex>

      <div className="position-feed">
        {!resource.data && !resource.error ? (
          Array.from({ length: 6 }, (_, index) => <FeedSkeleton key={index} />)
        ) : resource.data?.positions.length ? (
          resource.data.positions.map((position, index) => (
            <PositionRow
              key={position.address}
              position={position}
              emoji={POSITION_EMOJIS[index]}
            />
          ))
        ) : (
          <div className="feed-empty">
            <Text color={resource.error ? "red" : "gray"}>
              {resource.error?.message ?? "No indexed positions"}
            </Text>
          </div>
        )}
      </div>

      <Flex align="center" justify="between" className="feed-footer">
        <Text size="1" color="gray">Page {page}</Text>
        <Flex gap="2">
          <Button
            size="1"
            variant="soft"
            color="gray"
            disabled={page === 1}
            onClick={onPrevious}
          >
            Previous
          </Button>
          <Button
            size="1"
            variant="soft"
            disabled={!resource.data?.nextCursor}
            onClick={onNext}
          >
            Next
          </Button>
        </Flex>
      </Flex>
    </Card>
  );
}

function PositionRow({ position, emoji }: { position: PositionSummary; emoji: string }) {
  return (
    <div className="position-row">
      <span className="position-avatar" aria-hidden="true">{emoji}</span>
      <Text as="p" size="2" weight="medium" className="mono" title={position.address}>
        {shortAddress(position.address)}
      </Text>
      <PositionAmount label="Supply" value={position.supplyAssets} />
      <PositionAmount label="Borrow" value={position.borrowAssets} />
      <PositionAmount label="Collateral" value={position.collateral} />
    </div>
  );
}

function PositionAmount({ label, value }: { label: string; value: string }) {
  return (
    <div className="position-amount" title={value}>
      <Text as="p" size="1" color="gray">
        {label}
      </Text>
      <Text as="p" size="2" weight="medium" className="mono">
        {formatE18(value)}
      </Text>
    </div>
  );
}

function FeedSkeleton() {
  return (
    <div className="position-row">
      <Skeleton width="38px" height="38px" />
      <Skeleton width="130px" height="28px" />
      <Skeleton width="80px" height="28px" />
      <Skeleton width="80px" height="28px" />
      <Skeleton width="80px" height="28px" />
    </div>
  );
}

export default App;
