import React from 'react';
import { Link } from 'react-router-dom';

import PageHeader from 'Components/PageHeader';
import Widget from 'Components/Widget';
import { useTheme } from 'Containers/ThemeProvider';

const VulnDashboardPage = () => {
    const { isDarkMode } = useTheme();
    const entities = [
        'deployments',
        'images',
        'components',
        'cves',
        'clusters',
        'namespaces',
        'policies'
    ];

    const pageHeaderBGStyle = isDarkMode
        ? {}
        : {
              boxShadow: 'hsla(230, 75%, 63%, 0.62) 0px 5px 30px 0px',
              '--start': 'hsl(226, 70%, 60%)',
              '--end': 'hsl(226, 64%, 56%)'
          };

    return (
        <section className="flex flex-col relative min-h-full">
            <PageHeader
                classes="bg-gradient-horizontal z-10 sticky pin-t text-base-100"
                bgStyle={pageHeaderBGStyle}
                header="Vulnerability Management"
                subHeader="Dashboard"
            >
                <div className="flex flex-1 justify-end">
                    <p>Buttons go here</p>
                </div>
            </PageHeader>
            <div
                className="relative bg-gradient-diagonal p-6 xxxl:p-8"
                style={{ '--start': 'var(--base-200)', '--end': 'var(--primary-100)' }}
                id="capture-dashboard"
            >
                <div
                    className="grid grid-gap-6 xxxl:grid-gap-8 md:grid-auto-fit xxl:grid-auto-fit-wide md:grid-dense"
                    style={{ '--min-tile-height': '160px' }}
                >
                    {entities.map(entity => (
                        <Widget
                            key={entity}
                            className="s-2 overflow-hidde pdf-page"
                            header={entity}
                            headerComponents={
                                <Link
                                    to={`/main/vulnerability-management/${entity}`}
                                    className="no-underline"
                                >
                                    <button className="btn-sm btn-base" type="button">
                                        View All
                                    </button>
                                </Link>
                            }
                        >
                            <p>Widget contents go here.</p>
                        </Widget>
                    ))}
                </div>
            </div>
        </section>
    );
};
export default VulnDashboardPage;
